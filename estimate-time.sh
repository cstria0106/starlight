sudo ctr-starlight pull postgres:12-starlight && mkdir /tmp/test-pg-data && \
sudo ctr-starlight create \
        --mount type=bind,src=/tmp/test-pg-data,dst=/var/lib/postgresql/data,options=rbind:rw --env-file ./demo/config/all.env \
        --net-host \
        postgres:12-starlight \
        postgres:12-starlight \
    instance3 && \
sudo ctr task start instance3