import grpc from 'k6/net/grpc';
import { check, group, sleep } from 'k6';

export let options = {
    summaryTrendStats: ["min", "max", "avg", "p(90)", "p(95)", "p(99)"],
    stages: [
        { duration: '1s', target: 1 },
        { duration: '120s', target: 50 },
        { duration: '10s', target: 0 },
    ]
};

const client = new grpc.Client();
client.load(['.'], 'helloworld.proto');

const client2 = new grpc.Client();
client2.load(['.'], 'helloworld.proto');

export default () => {
    // Connect the gRPC client/s on first iteration
    if (__ITER == 0) {
        client.connect(`k8s:///${__ENV.GRPC_SERVER}:50051`, {
            plaintext: true
        });

        if (__ENV.GRPC_SERVER_2) {
            client2.connect(`k8s:///${__ENV.GRPC_SERVER_2}:50051`, {
                plaintext: true
            });
        }
    }

    group('demoGrpcGroup', function () {
        const data = { name: `ICaRUS client=1 VU=${__VU} Iter=${__ITER}` };
        const response = client.invoke('helloworld.Greeter/SayHello', data);

        check(response, {
            'status is OK': (r) => r && r.status === grpc.StatusOK,
        });
    });

    if (__ENV.GRPC_SERVER_2) {
        group('demoGrpcGroup2', function () {
            const data = { name: `ICaRUS client=2 VU=${__VU} Iter=${__ITER}` };
            const response = client2.invoke('helloworld.Greeter/SayHello', data);

            check(response, {
                'status2 is OK': (r) => r && r.status === grpc.StatusOK,
            });
        });
    }

    sleep(0.3);
};

export function teardown(data) {
    client.close();
}
