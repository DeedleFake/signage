FROM scratch
MAINTAINER DeedleFake <deedlefake@users.noreply.github.com>

EXPOSE 8080
VOLUME [ "/etc/ssl" ]

ADD "signage" "/signage"

ENTRYPOINT [ "/signage" ]
