## Report

SSE server for experiment events

:construction: Under development. Currently serving experiment.com's request status codes on the `hits` namespace. To subscribe (from a browser supporting SSEs):

```javascript
hits = new EventSource("https://experiment-report.herokuapp.com/subscribe/hits");
hits.onmessage = function(code) {
  console.log(code);
};
```
