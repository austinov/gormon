const period = 15*60*1000 // 15 minutes

var memData = [[new Date(), 0, 0, 0]];
var cpuData = [[new Date(), 0, 0, 0]];
var firstDeleted = false;
var eventSource;

var memGraph = new Dygraph(
  document.getElementById('graph-mem'),
    memData,
    {
      title: 'Memory Usage',
      titleHeight: 32,                                  
      // labels for hosts
      labels: [
          'Date',
          '109.234.37.213',
          '192.168.1.104',
          '192.168.1.231'
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
                  return humanBytes(y);
              },
              axisLabelFormatter: function(y) {
                  return humanBytes(y);
              },
              axisLabelWidth: 60
          },
      },
      legend: 'always',
        highlightCircleSize: 2,
        strokeWidth: 1,
        strokeBorderWidth: 1,

        highlightSeriesOpts: {
          strokeWidth: 3,
          strokeBorderWidth: 1,
          highlightCircleSize: 5
        },

      labelsUTC: false,
      visibility: [
          true,
          true,
          true
      ]
    }
);

var cpuGraph = new Dygraph(
  document.getElementById('graph-cpu'),
    cpuData,
    {
      title: 'CPU Usage',
      titleHeight: 32,                                  
      // labels for hosts
      labels: [
          'Date',
          '109.234.37.213',
          '192.168.1.104',
          '192.168.1.231'
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
                  return y + '%';
              },
              axisLabelFormatter: function(y) {
                  return y + '%';
              },
              axisLabelWidth: 60
          },
      },
      legend: 'always',
        highlightCircleSize: 2,
        strokeWidth: 1,
        strokeBorderWidth: 1,

        highlightSeriesOpts: {
          strokeWidth: 3,
          strokeBorderWidth: 1,
          highlightCircleSize: 5
        },

      labelsUTC: false,
      visibility: [
          true,
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
            // Memory usage
            var other2 = 0;
            var other3 = 0;
            if (firstDeleted && memData.length > 0) {
                var last = memData[memData.length-1];
                other2 = last[2];
                other3 = last[3];
            }
            memData.push([
                    new Date(stats.tstamp * 1000),
                    parseStat(stats.used_memory),
                    parseStat(other2),
                    parseStat(other3),
            ]);
            // CPU usage
            other2 = 0;
            other3 = 0;
            if (firstDeleted && cpuData.length > 0) {
                var last = cpuData[cpuData.length-1];
                other2 = last[2];
                other3 = last[3];
            }
            cpuData.push([
                    new Date(stats.tstamp * 1000),
                    parseStat(stats.used_cpu_perc, true),
                    parseStat(other2, true),
                    parseStat(other3, true),
            ]);
        }
        else if (stats.host === "192.168.1.104") {
            // Memory usage
            var other1 = 0;
            var other3 = 0;
            if (firstDeleted && memData.length > 0) {
                var last = memData[memData.length-1];
                other1 = last[1];
                other3 = last[3];
            }
            memData.push([
                    new Date(stats.tstamp * 1000),
                    parseStat(other1),
                    parseStat(stats.used_memory),
                    parseStat(other3)
            ]);
            // CPU usage
            other1 = 0;
            other3 = 0;
            if (firstDeleted && cpuData.length > 0) {
                var last = cpuData[cpuData.length-1];
                other1 = last[1];
                other3 = last[3];
            }
            cpuData.push([
                    new Date(stats.tstamp * 1000),
                    parseStat(other1, true),
                    parseStat(stats.used_cpu_perc, true),
                    parseStat(other3, true)
            ]);
        }
        else if (stats.host === "192.168.1.231") {
            // Memory usage
            var other1 = 0;
            var other2 = 0;
            if (firstDeleted && memData.length > 0) {
                var last = memData[memData.length-1];
                other1 = last[1];
                other2 = last[2];
            }
            memData.push([
                    new Date(stats.tstamp * 1000),
                    parseStat(other1),
                    parseStat(other2),
                    parseStat(stats.used_memory)
            ]);
            // CPU usage
            other1 = 0;
            other2 = 0;
            if (firstDeleted && cpuData.length > 0) {
                var last = cpuData[cpuData.length-1];
                other1 = last[1];
                other2 = last[2];
            }
            cpuData.push([
                    new Date(stats.tstamp * 1000),
                    parseStat(other1, true),
                    parseStat(other2, true),
                    parseStat(stats.used_cpu_perc, true)
            ]);
        }
        
        if (!firstDeleted) {
            memData.splice(0, 1);
            cpuData.splice(0, 1);
            firstDeleted = true;
        }
        
        var firstDate = memData[0][0];
        var lastDate = memData[memData.length-1][0];
        for (;(lastDate-firstDate)>period;) {
            memData.splice(0, 1);
            cpuData.splice(0, 1);
            firstDate = memData[0][0];
            lastDate = memData[memData.length-1][0];
        }
        memGraph.updateOptions( { 'file': memData } );
        cpuGraph.updateOptions( { 'file': cpuData } );
    };
    document.getElementById('start-btn').disabled = true;
    document.getElementById('stop-btn').disabled = false;
}

function parseStat(value, isFloat) {
    var v =  isFloat ? parseFloat(value) : parseInt(value);
    return isNaN(v) ? 0 : v;
}

function stop() {
    eventSource.close();
    document.getElementById('start-btn').disabled = false;
    document.getElementById('stop-btn').disabled = true;
}

function log(msg) {
    console.log(msg);
}

function humanBytes(bytes) {
    var unit = 1024;
    if (bytes < unit) {
        return bytes + 'B';
    }
    exp = parseInt(Math.log(bytes)  / Math.log(unit));
    pre = 'KMGTPE'.charAt(exp-1);
    return (bytes / Math.pow(unit, exp)).toFixed(2) + ' ' + pre;
}

start();
