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

export default () => {
    if (__ITER == 0) {
        client.connect(`k8s:///${__ENV.GRPC_SERVER}:50051`, {
            plaintext: true
        });
    }

    group('demoGrpcGroup', function () {
        const data = { name: 'ICaRUS ' + __VU + ' ' + __ITER };
        const response = client.invoke('helloworld.Greeter/SayHello', data);

        check(response, {
            'status is OK': (r) => r && r.status === grpc.StatusOK,
        });
    });

     sleep(1);
};

export function teardown(data) {
    client.close();
}
