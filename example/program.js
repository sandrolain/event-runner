addMetadata("foo", "bar");
setMetadata("foo", "bar2");
addMetadata("foo", "bar3");
setData({
  time: new Date().toISOString(),
  data: message.topic(),
});
setConfig("topic", `http://localhost:8989/${message.topic()}`);
setConfig("method", "POST");
