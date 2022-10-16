let x = [];
let y = [];

const fetchData = () => {
  axios
    .get(
      'https://api.coingecko.com/api/v3/coins/bitcoin/market_chart/range?vs_currency=usd&from=1638244800&to=1655524800'
    )
    .then(function (response) {
      // handle success
      x = response.data['prices'].map((price) => {
        return moment(price[0]).format();
      });
      y = response.data['prices'].map((price) => {
        return price[1];
      });
    })
    .then(() => {
      console.log(x);
      console.log(y);
      drawPlot();
    })
    .catch(function (error) {
      // handle error
      console.log(error);
    })
    .finally(function () {
      // always executed
    });
};

const drawPlot = () => {
  var data = [
    {
      x: x,
      y: ,

      type: 'scatter',
    },
  ];

  Plotly.newPlot('myDiv', data);
};

fetchData();
