const period = 15*60*1000 // 15 minutes

var data = [[new Date(), 0, 0]];
var firstDeleted = false;
var eventSource;

var g = new Dygraph(
  document.getElementById('graph-mem'),
    data,
    {
      // labels for hosts
      labels: [
          'Date',
          '109.234.37.213',
          '192.168.1.104'
      ],
      // Y-axis for stats by time
      vds: {
         drawPoints: true,
         pointSize: 2
      },
      vbox: {
         drawPoints: true,
         pointSize: 2
       },
      axes: {
          y: {
              valueFormatter: function(y) {
                  return y;
              },
              axisLabelFormatter: function(y) {
                  return y;
              },
              axisLabelWidth: 60
          },
      },
      legend: 'always',
      labelsUTC: false,
      visibility: [
          true,
          true
      ]
    }
);

function start() {
  if (!window.EventSource) {
    alert('This browser doesn\'t support EventSource.');
    return;
  }

  eventSource = new EventSource('/stat/');

  eventSource.onerror = function(e) {
    if (this.readyState == EventSource.CONNECTING) {
      log("Reconnecting...");
    } else {
      log("Connection error: " + this.readyState);
    }
  };

  eventSource.onmessage = function(e) {
    stats = JSON.parse(e.data);
    if (stats.host === "109.234.37.213") {
        var other = 0;
        if (firstDeleted && data.length > 0) {
          var last = data[data.length-1];
          other = last[2];
        }

        data.push([
            new Date(stats.tstamp * 1000),
            parseStat(stats.used_memory),
            parseStat(other),
        ]);
    }
    else if (stats.host === "192.168.1.104") {
        var other = 0;
        if (firstDeleted && data.length > 0) {
          var last = data[data.length-1];
          other = last[1];
        }
        data.push([
            new Date(stats.tstamp * 1000),
            parseStat(other),
            parseStat(stats.used_memory)
        ]);
    }

    if (!firstDeleted) {
        data.splice(0, 1);
        firstDeleted = true;
    }

    var firstDate = data[0][0];
    var lastDate = data[data.length-1][0];
    for (;(lastDate-firstDate)>period;) {
        data.splice(0, 1);
        firstDate = data[0][0];
        lastDate = data[data.length-1][0];
    }
    //console.log(data);

    g.updateOptions( { 'file': data } );
  };
}

function parseStat(value) {
  var v =  parseInt(value);
  return isNaN(v) ? 0 : v;
}

function stop() {
  eventSource.close();
}

function log(msg) {
  //logElem.innerHTML += msg + "<br>";
  console.log(msg);
}

start();
