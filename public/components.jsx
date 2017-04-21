// See for help: https://github.com/reactjs/redux/blob/master/examples/counter/src/index.js
class PodsContainer extends React.Component {
  static propTypes = {
    data: React.PropTypes.object.isRequired
  }

  // incrementIfOdd = () => {
  //   if (this.props.value % 2 !== 0) {
  //     this.props.onIncrement()
  //   }
  // }

  // incrementAsync = () => {
  //   setTimeout(this.props.onIncrement, 1000)
  // }

  render() {
    const { data } = this.props
    return (
      <p>Pods: {Object.keys(data.pods).length}</p>
      // <p>
      //   Clicked: {value} times
      //   {' '}
      //   <button onClick={onIncrement}>
      //     +
      //   </button>
      //   {' '}
      //   <button onClick={onDecrement}>
      //     -
      //   </button>
      //   {' '}
      //   <button onClick={this.incrementIfOdd}>
      //     Increment if odd
      //   </button>
      //   {' '}
      //   <button onClick={this.incrementAsync}>
      //     Increment async
      //   </button>
      // </p>
    )
  }
}
