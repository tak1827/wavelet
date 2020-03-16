FROM golang:1.13-alpine as build

# **************************************************
# Build Rheia
# **************************************************
RUN apk add --no-cache git

RUN mkdir /src
COPY go.mod /src/go.mod
COPY go.sum /src/go.sum
RUN (cd /src; go mod download)

ADD . /src
RUN (cd /src/cmd/wavelet; go build)
RUN (cd /src/cmd/client; go build connecter.go)

# **************************************************
# sshd
# **************************************************
FROM alpine:3.9.5

RUN apk add --no-cache openssh openrc

RUN rc-update add sshd && rc-status
RUN mkdir -p /run/openrc/ && touch /run/openrc/softlevel

RUN adduser -D ssh-user
# It will cause an ssh login error if user has no password
RUN echo "ssh-user:ssh-pass" | chpasswd

RUN apk add --no-cache sudo
RUN echo "ssh-user ALL=(ALL) NOPASSWD:ALL" >> /etc/sudoers

USER ssh-user
RUN mkdir $HOME/.ssh
RUN echo -e "-----BEGIN RSA PRIVATE KEY-----\n\
MIICWwIBAAKBgQCV/3YMILqCjc+wTfTj5yL5bBw3RynSNZsI/C7mJB+XN05Y6nFf\n\
vtQ9tKZsOmmDDBdBdO3vgnXYE9SakilBTaXLt4iBMRfHHmQLOEjgM4zWe9wsN+o+\n\
xERlxLqO2YQzqVPBF4j6iNaPnhTLHfdc76wRT44eu9cxhW1fvBFN6tVD/wIDAQAB\n\
AoGADHbsElnt5cdvCNGRlhaw+1UWiyloukUmGvsJ0qoRzfkUk6L3xOMQtAf9BGVe\n\
OzuX+yKCWzLWJByRwD6BGyIessx85LT7aOt5AAurTI8nbr6vzIcglZshk6GD3HtW\n\
PUt/92z5huh4gfwdY4cNgg061Ie8BU4YTgZIsjgqZoTaRvECQQDHckkemmUwFZVz\n\
wxhbtx/ezYWnmm+Zo82gQPqf3N3yeTDxDrZ8ZVWntf3QKe4wIWxj6Q/5bE32ZPVu\n\
rgVovjaVAkEAwIe7UlHb01QHjUXvNrH+JwtNlT4BkkTlFq3BxAIfYu8bz5BdVqWS\n\
HkL9mzimjs0Hg7vCmkq7AtL93RskcSFPQwJAdn7fFHCUt9nK/okl+q+1HrGqu+ZU\n\
ccrRS2nvvXUC5KG+QyoDn6qxFkBZOQkAe9YPyNhaZk1DkKVw8D4x9oY0dQJACnAR\n\
XsAjW8ylcCpHcr6pmbJKLXokGar0hIs9AwRecT2nU1RvSwBE8LEPELg8h2dkLecj\n\
ALxmcRvBOZh8HEoDAwJAPCYkfPYYDjfgdDP71H4PYaeBNh51L7/iJVIpRLt8cxLm\n\
j3BJTqi0p/8L2oqfteJfLGxljpSsM/flzMrxNOv/YA==\n\
-----END RSA PRIVATE KEY-----" > $HOME/.ssh/id_rsa
RUN chmod 600 $HOME/.ssh/id_rsa
RUN echo "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAAgQCV/3YMILqCjc+wTfTj5yL5bBw3RynSNZsI/C7mJB+XN05Y6nFfvtQ9tKZsOmmDDBdBdO3vgnXYE9SakilBTaXLt4iBMRfHHmQLOEjgM4zWe9wsN+o+xERlxLqO2YQzqVPBF4j6iNaPnhTLHfdc76wRT44eu9cxhW1fvBFN6tVD/w== ssh-user@example.com" >> $HOME/.ssh/authorized_keys

EXPOSE 22

# **************************************************
# supervisor
# **************************************************
USER root
RUN apk add --no-cache --update supervisor vim

RUN mkdir /home/ssh-user/log

COPY supervisord.conf /etc/supervisord.conf
# CMD service supervisord start

# **************************************************
# Rheia
# **************************************************

COPY --from=build /src/cmd/wavelet/ /home/ssh-user/
COPY --from=build /src/cmd/client/ /home/ssh-user/

EXPOSE 3000
EXPOSE 9000

# ENTRYPOINT ["/src/cmd/wavelet/wavelet", "-api.port", "9000", "-host", "3.20.255.76", "--wallet", "./wallet-tak.txt"]

CMD /etc/init.d/sshd start && tail -f /dev/null
