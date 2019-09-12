cc.Class({
    extends: cc.Component,

    properties: {
        //事件接收对象
        dataEventHandler:null,
        //游戏信息
        gameInfo:null,
        websocket: null,
    },

    //初始化网络事件
    initHandlers:function(){
        var self = this;
        this.websocket.addHandler("LoginReply",proto.LoginReply.deserializeBinary,function(data){
            console.log("LoginReply data in",data);
            //self.sendEvent('netmsg_disconnect',data);
        });

        this.websocket.addHandler("RegisterReply",proto.RegisterReply.deserializeBinary,function(data){
            console.log("RegisterReply data in",data);
            //self.sendEvent('netmsg_disconnect',data);
        });

        this.websocket.addHandler("ErrorMessage",proto.ErrorMessage.deserializeBinary,function(data){
            console.log("ErrorMessage data in",data);
            //self.sendEvent('netmsg_disconnect',data);
        });
        //网络断开
        /*cc.vv.Net.addHandler("disconnect",function(data){
            console.log("disconnect",data);
            self.sendEvent('netmsg_disconnect',data);
        });*/
    },
    
    //设置事件接收对像
    setEventHandle:function(handle){
        this.dataEventHandler = handle;
    },

    sendEvent(event, data){
        if(this.dataEventHandler){
            this.dataEventHandler.emit(event, data);
        }    
    },

    initSocket:function(){
        var websocket = require("./net/websocket");
        this.websocket = websocket
    },

    connectGameServer:function(callback){
        this.websocket.connect("ws://127.0.0.1:7459/ws");
        var self = this;
        var onConnectOK = function(){
            console.log("connectGameServer_onConnectOK");
            setTimeout(function(){
                let loginMessage = new proto.LoginRequest();
                loginMessage.setAccount("1111")
                loginMessage.setPassword("111")
                let loginMessageBytes = loginMessage.serializeBinary();
                let clientRequestMessageContent = new proto.google.protobuf.Any()
                clientRequestMessageContent.pack(loginMessageBytes,"LoginRequest")
                let clientRequestMessage = new proto.ClientRequest();
                clientRequestMessage.setComponentname("Login");
                clientRequestMessage.setMethodname("Login");
                clientRequestMessage.setMessagecontent(clientRequestMessageContent);
                let bytes = clientRequestMessage.serializeBinary();
                self.websocket.send_data(bytes);

                cc.log("序列化为字节:" + bytes);        
                let data = proto.ClientRequest.deserializeBinary(bytes);
                cc.log("反序列化为对象:" ,data);
                cc.log("反序列化any为对象:" ,data.getMessagecontent().unpack(proto.LoginRequest.deserializeBinary,"LoginRequest"));

                if(callback){
                    callback(true);
                }
            },2000);  
        };
        
        var onConnectFailed = function(){
            console.log("connectGameServer_failed!");
            if(callback){
                callback(false);
            }
        };

        this.websocket.initHandle(onConnectOK,onConnectFailed);
    }
});
