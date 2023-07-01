// colors
var blue = '#0d66ff';
var blueLight = '#0e97b5';
var purple = '#6C5DD3';
var white = '#ffffff';
// var blueOpacity = '#e6efff';
// var blueLight = '#50B5FF';
var pink = '#FFB7F5';
// var orangeOpacity = '#fff5ed';
var yellow = '#FFCE73';
var green = '#7FBA7A';
var red = '#FF754C';
// var greenOpacity = '#ecfbf5';
var gray = '#808191';
var grayOpacity = '#f2f2f2';
// var grayLight = '#E2E2EA';
var borderColor = "#E4E4E4";
// var text = "#171725";

// charts
Apex.chart = {
  fontFamily: 'Inter, sans-serif',
  fontSize: 13,
  fontWeight: 500,
  foreColor: gray
};


 
 
 

// chart users blue color
(function () {
  var options = {
    labels: ['Jan', 'Feb', 'Mar', 'Apr', 'May'],
    series: [{
      name: '',
      data: [18, 34, 44, 58, 38]
    }],
    colors: [red],
    chart: {
      height: '50%',
      type: 'bar',
      toolbar: {
        show: false
      },
      sparkline: {
        enabled: true
      }
    },
    grid: {
      borderColor: borderColor,
      strokeDashArray: 0,
      xaxis: {
        lines: {
          show: true
        }
      },
      yaxis: {
        lines: {
          show: false
        }
      }
    },
    stroke: {
      curve: 'smooth'
    },
    dataLabels: {
      enabled: false
    },
    plotOptions: {
      bar: {
        columnWidth: '80%'
      }
    },
    states: {
      hover: {
        filter: {
          type: 'darken',
          value: 0.9
        }
      }
    },
    legend: {
      show: false
    },
    tooltip: {
      shared: true
    },
    xaxis: {
      axisBorder: {
        show: false,
        color: borderColor
      },
      axisTicks: {
        show: false
      },
      tooltip: {
        enabled: false
      }
    }
  };

  var chart = document.querySelector('#chart-users-blue');
  if (chart != null) {
    new ApexCharts(chart, options).render();
  }
})();

// chart users blue color
(function () {
  var options = {
    labels: ['Jan', 'Feb', 'Mar', 'Apr', 'May'],
    series: [{
      name: '',
      data: [35, 66, 34, 56, 18]
    }],
    colors: [pink],
    chart: {
      height: '50%',
      type: 'bar',
      toolbar: {
        show: false
      },
      sparkline: {
        enabled: true
      }
    },
    grid: {
      borderColor: borderColor,
      strokeDashArray: 0,
      xaxis: {
        lines: {
          show: true
        }
      },
      yaxis: {
        lines: {
          show: false
        }
      }
    },
    stroke: {
      curve: 'smooth'
    },
    dataLabels: {
      enabled: false
    },
    plotOptions: {
      bar: {
        columnWidth: '80%'
      }
    },
    states: {
      hover: {
        filter: {
          type: 'darken',
          value: 0.9
        }
      }
    },
    legend: {
      show: false
    },
    tooltip: {
      shared: true
    },
    xaxis: {
      axisBorder: {
        show: false,
        color: borderColor
      },
      axisTicks: {
        show: false
      },
      tooltip: {
        enabled: false
      }
    }
  };

  var chart = document.querySelector('#chart-users-blue1');
  if (chart != null) {
    new ApexCharts(chart, options).render();
  }
})();

// chart users blue color
(function () {
  var options = {
    labels: ['Jan', 'Feb', 'Mar', 'Apr', 'May'],
    series: [{
      name: '',
      data: [35, 66, 34, 56, 18]
    }],
    colors: [blue],
    chart: {
      height: '50%',
      type: 'bar',
      toolbar: {
        show: false
      },
      sparkline: {
        enabled: true
      }
    },
    grid: {
      borderColor: borderColor,
      strokeDashArray: 0,
      xaxis: {
        lines: {
          show: true
        }
      },
      yaxis: {
        lines: {
          show: false
        }
      }
    },
    stroke: {
      curve: 'smooth'
    },
    dataLabels: {
      enabled: false
    },
    plotOptions: {
      bar: {
        columnWidth: '80%'
      }
    },
    states: {
      hover: {
        filter: {
          type: 'darken',
          value: 0.9
        }
      }
    },
    legend: {
      show: false
    },
    tooltip: {
      shared: true
    },
    xaxis: {
      axisBorder: {
        show: false,
        color: borderColor
      },
      axisTicks: {
        show: false
      },
      tooltip: {
        enabled: false
      }
    }
  };

  var chart = document.querySelector('#chart-users-blue2');
  if (chart != null) {
    new ApexCharts(chart, options).render();
  }
})();

// chart users blue color
(function () {
  var options = {
    labels: ['Jan', 'Feb', 'Mar', 'Apr', 'May', 'Feb', 'Mar', 'Apr', 'May'],
    series: [{
      name: '',
      data: [35, 66, 34, 56, 98, 66, 34, 56, 18]
    }],
    colors: [pink],
    chart: {
      height: '50%',
      type: 'bar',
      toolbar: {
        show: false
      },
      sparkline: {
        enabled: true
      }
    },
    grid: {
      borderColor: borderColor,
      strokeDashArray: 0,
      xaxis: {
        lines: {
          show: true
        }
      },
      yaxis: {
        lines: {
          show: false
        }
      }
    },
    stroke: {
      curve: 'smooth'
    },
    dataLabels: {
      enabled: false
    },
    plotOptions: {
      bar: {
        columnWidth: '80%'
      }
    },
    states: {
      hover: {
        filter: {
          type: 'darken',
          value: 0.9
        }
      }
    },
    legend: {
      show: false
    },
    tooltip: {
      shared: true
    },
    xaxis: {
      axisBorder: {
        show: false,
        color: borderColor
      },
      axisTicks: {
        show: false
      },
      tooltip: {
        enabled: false
      }
    }
  };

  var chart = document.querySelector('#chart-users-blue3');
  if (chart != null) {
    new ApexCharts(chart, options).render();
  }
})();
 
 
 

// chart report
(function () {
  var options = {
    labels: ['Oct', 'Nov', 'Dec', 'Jan', 'Feb'],
    series: [{
      name: 'New users',
      data: [70, 25, 44, 37, 28]
    }, {
      name: 'Users',
      data: [40, 16, 38, 30, 25]
    }],
    colors: [purple, blue],
    chart: {
      height: '100%',
      type: 'bar',
      toolbar: {
        show: false
      }
    },
    grid: {
      borderColor: borderColor,
      strokeDashArray: 0,
      xaxis: {
        lines: {
          show: false
        }
      },
      yaxis: {
        lines: {
          show: false
        }
      },
      padding: {
        top: 0,
        left: 10,
        right: 0,
        bottom: 0
      }
    },
    states: {
      hover: {
        filter: {
          type: 'darken',
          value: 0.9
        }
      }
    },
    stroke: {
      curve: 'smooth'
    },
    dataLabels: {
      enabled: false
    },
    plotOptions: {
      bar: {
        columnWidth: '60%'
      }
    },
    legend: {
      show: false
    },
    tooltip: {
      x: {
        show: false
      },
      shared: true
    },
    xaxis: {
      axisBorder: {
        show: false
      },
      axisTicks: {
        show: false
      },
      tooltip: {
        enabled: false
      }
    },
    yaxis: {
      axisBorder: {
        color: borderColor
      },
      axisTicks: {
        show: false
      },
      tooltip: {
        enabled: false
      }
    }
  };

  var chart = document.querySelector('#chart-report');
  if (chart != null) {
    new ApexCharts(chart, options).render();
  }
})();

 // chart earning
(function () {
    var options = {
      series: [80],
      chart: {
	      type: 'radialBar',
	      height: '150%',
    },
    plotOptions: {
      radialBar: {
        startAngle: -90,
        endAngle: 90,
        track: {
          background: borderColor,
          strokeWidth: '100%',
          margin: 0, // margin is in pixels
          
        },
        dataLabels: {
          name: {
            show: false
          },
          value: {
            offsetY: -2,
            fontSize: '32px',
            color: '#000',
            fontWeight: '700'
          }
        }
      }
    },
    grid: {
      padding: {
        top: -10
      }
    },
    fill: {
      type: 'gradient',
      gradient: {
        shade: 'light',
        shadeIntensity: 0.4,
        inverseColors: false,
        opacityFrom: 1,
        opacityTo: 1,
        stops: [0, 50, 53, 91]
      },
    },
  };

  var chart = document.querySelector('#chart-half-report');
  if (chart != null) {
    new ApexCharts(chart, options).render();
  }
})();


// chart earnings by item
(function () {
  var options = {
      series: [44, 55, 67],
      colors: [blue, pink, red],
      chart: {
      height: 350,
      type: 'radialBar',
    },
    plotOptions: {
      radialBar: {
        dataLabels: {
          name: {
            fontSize: '14px',
            color: '#FFF',
            show: false,
          },
          value: {
          	// offsetY: -1,
            fontSize: '16px',
            fontSize: '44px',
            color: '#000',
            fontWeight: '700'
          },
          total: {
            show: true,
            label: 'Total',
            formatter: function (w) {
              // By default this function returns the average of all series. The below is just an example to show the use of custom formatter function
              return 249
            }
          }
        }
      }
    },
    labels: ['Apples', 'Oranges', 'Bananas'],
    };

  var chart = document.querySelector('#chart-multipleitem');
  if (chart != null) {
    new ApexCharts(chart, options).render();
  }
})();

// chart earnings by item
(function () {
  var options = {
    series: [25, 37, 38],
    colors: [blue, pink, red],
    chart: {
      height: '140%',
      type: 'donut'
    },
    plotOptions: {
      pie: {
        donut: {
          size: '60%',
          polygons: {
            strokeWidth: 0
          }
        },
        expandOnClick: true
      }
    },
    dataLabels: {
      enabled: false
    },
    states: {
      hover: {
        filter: {
          type: 'darken',
          value: 0.9
        }
      }
    },
    legend: {
      show: false
    },
    tooltip: {
      enabled: true
    }
  };

  var chart = document.querySelector('#chart-earnings-by-item');
  if (chart != null) {
    new ApexCharts(chart, options).render();
  }
})();


// chart users blue color
(function () {
	var options = {
    labels: ['Jan', 'Feb', 'Mar', 'Apr', 'May'],
    series: [{
      name: '',
      data: [35, 66, 34, 56, 18]
    },{
      name: '',
      data: [12, 34, 12, 11, 7]
    }],
    colors: [blue,pink],
      chart: {
      type: 'bar',
      height: '100%',
      stacked: true,
      toolbar: {
        show: false
      },


      
    },
    responsive: [{
      breakpoint: 480,
      options: {
        legend: {
          position: 'bottom',
          offsetX: -10,
          offsetY: 0
        }
      }
    }],
    plotOptions: {
      bar: {
        horizontal: false,
      },
    },
    
    legend: {
      show: false
    },
    fill: {
      opacity: 1
    }
  };
   

  var chart = document.querySelector('#chart-usersMultiple');
  if (chart != null) {
    new ApexCharts(chart, options).render();
  }
})();

// chart users blue color
(function () {
	var options = {
    labels: ['Jan', 'Feb', 'Mar', 'Apr', 'May'],
    series: [{
      name: '',
      data: [35, 66, 34, 56, 18]
    },{
      name: '',
      data: [12, 34, 12, 11, 7]
    }],
    colors: [purple,pink],
      chart: {
      type: 'bar',
      height: '100%',
      stacked: true,
      toolbar: {
        show: false
      },
      sparkline: {
        enabled: true
      }      
      
    },
    responsive: [{
      breakpoint: 480,
      options: {
        legend: {
          position: 'bottom',
          offsetX: -10,
          offsetY: 0
        }
      }
    }],
    plotOptions: {
      bar: {
        horizontal: false,
      },
    },
    
    legend: {
      show: false
    },
    fill: {
      opacity: 1
    }
  };
   

  var chart = document.querySelector('#chart-usersMultipleX');
  if (chart != null) {
    new ApexCharts(chart, options).render();
  }
})();

// chart users blue color
(function () {
	var options = {
	  series: [{
	  data: [{
	      x: new Date(1538778600000),
	      y: [6629.81, 6650.5, 6623.04, 6633.33]
	    },
	    {
	      x: new Date(1538780400000),
	      y: [6632.01, 6643.59, 6620, 6630.11]
	    },
	    {
	      x: new Date(1538782200000),
	      y: [6630.71, 6648.95, 6623.34, 6635.65]
	    },
	    {
	      x: new Date(1538784000000),
	      y: [6635.65, 6651, 6629.67, 6638.24]
	    },
	    {
	      x: new Date(1538785800000),
	      y: [6638.24, 6640, 6620, 6624.47]
	    },
	    {
	      x: new Date(1538787600000),
	      y: [6624.53, 6636.03, 6621.68, 6624.31]
	    },
	    {
	      x: new Date(1538789400000),
	      y: [6624.61, 6632.2, 6617, 6626.02]
	    },
	    {
	      x: new Date(1538791200000),
	      y: [6627, 6627.62, 6584.22, 6603.02]
	    },
	    {
	      x: new Date(1538793000000),
	      y: [6605, 6608.03, 6598.95, 6604.01]
	    },
	    {
	      x: new Date(1538794800000),
	      y: [6604.5, 6614.4, 6602.26, 6608.02]
	    },
	    {
	      x: new Date(1538796600000),
	      y: [6608.02, 6610.68, 6601.99, 6608.91]
	    },
	    {
	      x: new Date(1538798400000),
	      y: [6608.91, 6618.99, 6608.01, 6612]
	    },

	  ]
	}],

	  chart: {
	  type: 'candlestick',
	  height: 350,
	  colors: [purple,pink],
	  toolbar: {
        show: false
      },
	    
	},
	
	xaxis: {
	  type: 'datetime'
	},
	yaxis: {
	  tooltip: {
	    enabled: false
	  }
	}
	};
   

  var chart = document.querySelector('#chart-candlestick');
  if (chart != null) {
    new ApexCharts(chart, options).render();
  }
})();




// chart users blue color
(function () {
	var options = {
      series: [{
      name: 'series1',
      data: [31, 40, 28, 51, 42, 109, 100]
    }],
    colors: [white],
      chart: {
      height: '50%',
      type: 'line',
      toolbar: {
        show: false
      },
      sparkline: {
        enabled: true
      }
    },
    dataLabels: {
      enabled: false
    },
    stroke: {
      width: 3,
      curve: 'smooth'
    },
    };
   

  var chart = document.querySelector('#chart-facebook');
  if (chart != null) {
    new ApexCharts(chart, options).render();
  }
})();

// chart users blue color
(function () {
	var options = {
      series: [{
      name: 'series1',
      data: [11, 32, 45, 32, 34, 52, 41]
    }],
    colors: [white],
      chart: {
      height: '50%',
      type: 'line',
      toolbar: {
        show: false
      },
      sparkline: {
        enabled: true
      }
    },
    dataLabels: {
      enabled: false
    },
    stroke: {
      width: 3,
      curve: 'smooth'
    },
    };
   

  var chart = document.querySelector('#chart-twiiter');
  if (chart != null) {
    new ApexCharts(chart, options).render();
  }
})();

// chart users blue color
(function () {
	var options = {
      series: [{
      name: 'series1',
      data: [31, 40, 28, 51, 42, 109, 100]
    }],
    colors: [white],
      chart: {
      height: '50%',
      type: 'line',
      toolbar: {
        show: false
      },
      sparkline: {
        enabled: true
      }
    },
    dataLabels: {
      enabled: false
    },
    stroke: {
      width: 3,
      curve: 'smooth'
    },
    };
   

  var chart = document.querySelector('#chart-instagram');
  if (chart != null) {
    new ApexCharts(chart, options).render();
  }
})();


// chart users blue color
(function () {
	var options = {
      series: [{
      name: 'series1',
      data: [31, 40, 28, 51, 42, 109, 100]
    }],
    colors: [blue],
      chart: {
      height: '50%',
      type: 'line',
      toolbar: {
        show: false
      },
      sparkline: {
        enabled: true
      }
    },
    dataLabels: {
      enabled: false
    },
    stroke: {
      width: 3,
      curve: 'smooth'
    },
   
    };
   

  var chart = document.querySelector('#chart-revinue');
  if (chart != null) {
    new ApexCharts(chart, options).render();
  }
})();

// chart users blue color
(function () {
	var options = {
      series: [{
      name: 'series1',
      data: [31, 40, 28, 51, 42, 109, 100]
    }],
    colors: [blue],
      chart: {
      height: '50%',
      type: 'line',
      toolbar: {
        show: false
      },
      sparkline: {
        enabled: true
      }
    },
    dataLabels: {
      enabled: false
    },
    stroke: {
      width: 3,
      curve: 'smooth'
    },
   
    };
   

  var chart = document.querySelector('#chart-revinuee');
  if (chart != null) {
    new ApexCharts(chart, options).render();
  }
})();

// chart users blue color
(function () {
	var options = {
      series: [{
      name: 'series1',
      data: [11, 32, 45, 32, 34, 52, 41]
    }],
    colors: [white],
      chart: {
      height: '60%',
      type: 'area',
      toolbar: {
        show: false
      },
      sparkline: {
        enabled: true
      }
    },
    dataLabels: {
      enabled: false
    },
    stroke: {
      width: 3,
      curve: 'smooth'
    },
   
    };
   

  var chart = document.querySelector('#chart-check');
  if (chart != null) {
    new ApexCharts(chart, options).render();
  }
})();


// chart users blue color
(function () {
	var options = {
    labels: ['Jan', 'Feb', 'Mar', 'Apr', 'May' ,'Jun', 'Jul', 'Aug', 'Sep', 'Oct' , 'Nov', 'Dec'],
    series: [{
      name: '',
      data: [35, 66, 34, 56, 18 ,35, 66, 34, 56, 18 , 56, 18]
    },{
      name: '',
      data: [12, 34, 12, 11, 7 ,12, 34, 12, 11, 7 , 11, 7]
    }],
    colors: [blue,pink],
      chart: {
      type: 'bar',
      height: '250%',
      stacked: true,
      toolbar: {
        show: false
      },
    },
    responsive: [{
      breakpoint: 480,
      options: {
        legend: {
          position: 'bottom',
          offsetX: -10,
          offsetY: 0
        }
      }
    }],
    plotOptions: {
      bar: {
        horizontal: false,
      },
    },
    
    legend: {
      show: false
    },
    fill: {
      opacity: 1
    },
  };
   

  var chart = document.querySelector('#chart-usersMultiplee');
  if (chart != null) {
    new ApexCharts(chart, options).render();
  }
})();

// chart users blue color
(function () {
  var options = {
    labels: ['Jan', 'Feb', 'Mar', 'Apr', 'May'],
    series: [{
      name: '',
      data: [6, 14, 12, 11, 18 ]
    },{
      name: '',
      data: [12, 7, 10, 11, 7]
    }],
    colors: [blue,pink],
      chart: {
      type: 'bar',
      height: '140%',
      stacked: true,
      toolbar: {
        show: false
      },
    },
    responsive: [{
      breakpoint: 480,
      options: {
        legend: {
          position: 'bottom',
          offsetX: -10,
          offsetY: 0
        }
      }
    }],
    plotOptions: {
      bar: {
        horizontal: false,
      },
    },
    
    legend: {
      show: false
    },
    fill: {
      opacity: 1
    },
  };
   

  var chart = document.querySelector('#chart-usersMultiplee2');
  if (chart != null) {
    new ApexCharts(chart, options).render();
  }
})();

// chart users blue color
(function () {
	 var options = {
      series: [{
      name: 'Marine Sprite',
      data: [76]
    }, {
      name: 'Striking Calf',
      data: [34]
    }],
    colors: [blue,borderColor],
      chart: {
      type: 'bar',
      height: 10,
      stacked: true,
      stackType: '100%',
      toolbar: {
        show: false
      },
      sparkline: {
        enabled: true
      }
    },
    plotOptions: {
      bar: {
        horizontal: true,
      },
    },
    
    fill: {
      opacity: 1
    
    },
    };

   

  var chart = document.querySelector('#chart-hor-bar');
  if (chart != null) {
    new ApexCharts(chart, options).render();
  }
})();

// chart users blue color
(function () {
	 var options = {
      series: [{
      name: 'Marine Sprite',
      data: [76]
    }, {
      name: 'Striking Calf',
      data: [34]
    }],
    colors: [red,borderColor],
      chart: {
      type: 'bar',
      height: 10,
      stacked: true,
      stackType: '100%',
      toolbar: {
        show: false
      },
      sparkline: {
        enabled: true
      }
    },
    plotOptions: {
      bar: {
        horizontal: true,
      },
    },
    
    fill: {
      opacity: 1
    
    },
    };

   

  var chart = document.querySelector('#chart-hor-bar2');
  if (chart != null) {
    new ApexCharts(chart, options).render();
  }
})();


// chart users blue color
(function () {
	var options = {
      series: [{
      name: 'Website Blog',
      type: 'column',
      data: [440, 505, 671, 543, 427]
    }, {
      name: 'Social Media',
      type: 'line',
      data: [323, 432, 601, 514, 343]
    }],
      chart: {
      height: 350,
      type: 'line',
      toolbar: {
        show: false
      }
    },
    stroke: {
      width: [0, 2]
    },
     
    dataLabels: {
      enabled: false,
      enabledOnSeries: [1]
    },
    labels: ['jan','Feb','Mar','Apr','May'],
    legend: {
      show: false
    },
    };

   

  var chart = document.querySelector('#chart-multyple-bar');
  if (chart != null) {
    new ApexCharts(chart, options).render();
  }
})();




// chart earnings by item
(function () {
  var options = {
      series: [67],
      colors: [red],
      chart: {
      height: '70%',
      type: 'radialBar',
    },
    plotOptions: {
      radialBar: {
        dataLabels: {
          name: {
            fontSize: '14px',
            color: '#FFF',
            show: false,
          },
          value: {
          	offsetY: -0.2,
            fontSize: '14px', 
            color: '#000',
            fontWeight: '700',
            marginTop:'5px'
          },
          total: {
            show: true,
            label: 'Total',
            formatter: function (w) {
              // By default this function returns the average of all series. The below is just an example to show the use of custom formatter function
              return 67+'%';
            }
          }
        }
      }
    },
    // labels: ['Apples', 'Oranges', 'Bananas'],
    };

  var chart = document.querySelector('#chart-round-progress');
  if (chart != null) {
    new ApexCharts(chart, options).render();
  }
})();



// chart earnings by item
(function () {
   var options = {
      series: [{
      name: 'TEAM A',
      type: 'column',
      data: [44, 11, 22, 27, 13, 22, 37]
    }, {
      name: 'TEAM B',
      type: 'area',
      data: [44, 55, 41, 67, 22, 43, 21]
    }, {
      name: 'TEAM C',
      type: 'line',
      data: [30, 25, 36, 30, 45, 35, 64]
    }],
      chart: {
      height: '150%',
      type: 'line',
      stacked: false,
      toolbar: {
        show: false
      },
      sparkline: {
        enabled: true
      }
    },
    stroke: {
      width: [0, 2, 2],
      curve: 'smooth'
    },
    plotOptions: {
      bar: {
        columnWidth: '50%'
      }
    },
    
    fill: {
      opacity: [0.85, 0.25, 1],
      gradient: {
        inverseColors: false,
        shade: 'light',
        type: "vertical",
        opacityFrom: 0.85,
        opacityTo: 0.55,
        stops: [0, 100, 100, 100]
      }
    },
    labels: ['01/01/2003', '02/01/2003', '03/01/2003', '04/01/2003', '05/01/2003', '06/01/2003', '07/01/2003'],
    markers: {
      size: 0
    },
    xaxis: {
      type: 'datetime'
    },
    yaxis: {
      title: {
        text: 'Points',
      },
      min: 0
    },
    tooltip: {
      shared: true,
      intersect: false,
      y: {
        formatter: function (y) {
          if (typeof y !== "undefined") {
            return y.toFixed(0) + " points";
          }
          return y;
    
        }
      }
    }
    };

  var chart = document.querySelector('#chart-center');
  if (chart != null) {
    new ApexCharts(chart, options).render();
  }
})();


// chart earnings by item
(function () {
   

    var options = {
          series: [14, 23, 21, 17, 15, 10, 12, 17, 21],
          chart: {
          type: 'polarArea',
          width: '100%',
          height:'250%',
        },
        stroke: {
          colors: ['#fff']
        },
        fill: {
          opacity: 0.8
        },
        responsive: [{
          breakpoint: 480,
          options: {
            chart: {
              width: '100%',
            },
            legend: {
              position: 'bottom'
            }
          }
        }],
        legend: {
	      show: false
	    },
        };

  var chart = document.querySelector('#chart-round-center');
  if (chart != null) {
    new ApexCharts(chart, options).render();
  }
})();



// chart earnings by item
(function () {
   
  var options = {
    series: [{
    name: 'Net Profit',
    data: [44, 55, 57, 56, 61, 58, 63, 60, 66]
  }, {
    name: 'Revenue',
    data: [76, 85, 101, 98, 87, 105, 91, 114, 94]
  }, {
    name: 'Free Cash Flow',
    data: [35, 41, 36, 26, 45, 48, 52, 53, 41]
  }],
    chart: {
    type: 'bar',
    height: 350,
    stacked: true,
    toolbar: {
        show: false
      },
  },
  plotOptions: {
    bar: {
      horizontal: false,
      columnWidth: '45%',
      endingShape: 'rounded'
    },
  },
  dataLabels: {
    enabled: false
  },
  stroke: {
    show: true,
    width: 2,
    // colors: ['transparent']
  },
  xaxis: {
    categories: ['Feb', 'Mar', 'Apr', 'May', 'Jun', 'Jul', 'Aug', 'Sep', 'Oct'],
  },
  
  fill: {
    opacity: 1
  },
   legend: {
      show: false
    },
  };

  var chart = document.querySelector('#column-round-chart');
  if (chart != null) {
    new ApexCharts(chart, options).render();
  }
})();

