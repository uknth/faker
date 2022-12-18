FROM ubuntu:22.04
WORKDIR /faker

ADD bin/faker /faker
ADD fake.yaml /faker
ADD scripts/run.sh /faker

ENV PATH="${PATH}:/faker"

CMD ["bash", "run.sh"]