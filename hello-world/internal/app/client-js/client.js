var PROTO_PATH = '../../../api/proto/hello.proto';

var grpc = require('@grpc/grpc-js');
var grpc_xds = require('@grpc/grpc-js-xds');
grpc_xds.register();

var protoLoader = require('@grpc/proto-loader');
var packageDefinition = protoLoader.loadSync(
    PROTO_PATH,
    {keepCase: true,
        longs: String,
        enums: String,
        defaults: true,
        oneofs: true
    });
var hello = grpc.loadPackageDefinition(packageDefinition).hello;

function main() {
    target = 'xds:///hello-server-listener';

    var client = new hello.Hello(target, grpc.credentials.createInsecure());

    client.SayHello({name: "gRPC Proxyless LB-js"}, function(err, response) {
        if (err) throw err;
        console.log('Greeting:', response.message);
        client.close();
    });
}

main();