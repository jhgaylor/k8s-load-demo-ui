const rootEl = document.getElementById('root')
const store = Redux.createStore(Reducers.get())

Websocket.setup(store)

const render = () => ReactDOM.render(
  <PodsContainer
    data={store.getState()}
  />,
  rootEl
)

render()
store.subscribe(render)
