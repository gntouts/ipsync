FROM scratch

COPY ./bin/ipsync /usr/bin/

CMD ["ipsync"]