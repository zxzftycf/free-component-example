cc.Class({
    extends: cc.Component,

    properties: {
        label: {
            default: null,
            type: cc.Label
        },
        // defaults, set visually when attaching this script to the Canvas
        text: 'Hello, World!',
        gameNetMgr : null
    },

    // use this for initialization
    onLoad: function () {
        this.label.string = this.text;

        let registerMessage = new proto.RegisterRequest();
        registerMessage.setAccount("1111")
        registerMessage.setPassword("111")
        let registerMessageBytes = registerMessage.serializeBinary();
        let clientRequestMessageContent = new proto.google.protobuf.Any()
        clientRequestMessageContent.pack(registerMessageBytes,"RegisterRequest")
        let clientRequestMessage = new proto.ClientRequest();
        clientRequestMessage.setComponentname("Register");
        clientRequestMessage.setMethodname("Register");
        clientRequestMessage.setMessagecontent(clientRequestMessageContent);
        cc.log("register data:" ,clientRequestMessage);
        let bytes = clientRequestMessage.serializeBinary();
        cc.log("序列化为字节:" + bytes);   
        var http = require("./net/http");
        var httper = new http()
        var gameNetMgr = require("./gameNetMgr");
        this.gameNetMgr = new gameNetMgr();
        this.gameNetMgr.initSocket()
        this.gameNetMgr.initHandlers();
        var self = this;
        httper.sendWithUrl("",bytes, function(reply){
            let data = proto.ClientReply.deserializeBinary(reply);
            let messageName = data.getMessagename()
            self.gameNetMgr.websocket.trigger(messageName,data)

            self.gameNetMgr.connectGameServer();
        })
    },

    // called every frame
    update: function (dt) {

    },
});
