addMetadata("foo", "bar");
setMetadata("foo", "bar2");
addMetadata("foo", "bar3");
setConfig("topic", `http://localhost:8989/${message.topic()}?foo=bar&foo=baz`);
setConfig("method", "POST");

const foo = cache.get("foo") ?? 0;
const bar = cache.get("bar") ?? "";

const time = new Date().toISOString();

setData({
  time,
  data: message.topic(),
  foo,
  bar,
});

cache.set("foo", foo + 1);
cache.set("bar", `bar: ${time}`);
