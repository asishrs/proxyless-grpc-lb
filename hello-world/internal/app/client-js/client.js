var PROTO_PATH = './api/proto/hello.proto';

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


function sleep(ms) {
    return new Promise(resolve => setTimeout(resolve, ms));
}

function main() {
    target = 'xds:///hello-server-listener';
    console.info("starting client")
    var client = new hello.Hello(target, grpc.credentials.createInsecure());
    console.info("created client")
    console.info(client)

        client.SayHello({name: "gRPC Proxyless LB-js"}, function(err, response) {
            if (err){
                console.error(err)
                client.close()
                client = new hello.Hello(target, grpc.credentials.createInsecure());
            }

            console.log('Greeting:', response.message);
        });
}

main();