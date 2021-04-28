from common import Runner
from testcases import *

if __name__ == '__main__':
    t = TestRedis("6.0", "5.0")
    t.rounds = 2
    r = Runner()
    history_temp = []

    r.service.reset_latency_bandwidth()
    print("Hello! This is Starlight Stage.")

    pool_starlight = []
    pool_vanilla = []
    pool_estargz = []
    pool_wget = []

    for i in range(len(t.rtt)):
        print("RTT:%d" % t.rtt[i])
        step_starlight = []
        step_vanilla = []
        step_estargz = []
        step_wget = []

        # estargz
        r.service.reset_container_service()
        r.service.start_grpc_estargz()

        n = r.sync_pull_estargz(t, 0, False)

        r.service.set_latency_bandwidth(t.rtt[i])
        r.test_estargz(t, history_temp, n, False)
        r.service.reset_latency_bandwidth()

        r.service.kill_estargz()

        for k in range(t.rounds):
            r.service.reset_container_service()
            r.service.start_grpc_estargz()

            n = r.sync_pull_estargz(t, 0, False)

            r.service.set_latency_bandwidth(t.rtt[i])
            r.test_estargz(t, history_temp, n, False)
            r.service.reset_latency_bandwidth()

            r.service.kill_estargz()

        # starlight
        r.service.reset_container_service()
        r.service.start_grpc_starlight()

        n = r.sync_pull_starlight(t, 0, False)

        r.service.set_latency_bandwidth(t.rtt[i])
        r.test_starlight(t, history_temp, n, False)
        r.service.reset_latency_bandwidth()

        r.service.kill_starlight()

        for k in range(t.rounds):
            r.service.reset_container_service()
            r.service.start_grpc_starlight()

            n = r.sync_pull_starlight(t, 0, False)

            r.service.set_latency_bandwidth(t.rtt[i])
            r.test_starlight(t, history_temp, n, False)
            r.service.reset_latency_bandwidth()

            r.service.kill_starlight()

        # vanilla
        r.service.reset_container_service()

        n = r.sync_pull_vanilla(t, 0, False)

        r.service.set_latency_bandwidth(t.rtt[i])
        r.test_vanilla(t, history_temp, n, False)
        r.service.reset_latency_bandwidth()

        for k in range(t.rounds):
            r.service.reset_container_service()

            n = r.sync_pull_vanilla(t, 0, False)

            r.service.set_latency_bandwidth(t.rtt[i])
            r.test_vanilla(t, history_temp, n, False)
            r.service.reset_latency_bandwidth()

        # wget
        r.service.set_latency_bandwidth(t.rtt[i])
        r.test_wget(t, history_temp)
        r.service.reset_latency_bandwidth()

        for k in range(t.rounds):
            r.service.set_latency_bandwidth(t.rtt[i])
            r.test_wget(t, history_temp)
            r.service.reset_latency_bandwidth()

        pool_starlight.append(step_starlight)
        pool_vanilla.append(step_vanilla)
        pool_estargz.append(step_estargz)
        pool_wget.append(step_wget)

        t.save_results(pool_starlight, pool_vanilla, pool_estargz, pool_wget, i + 1)
