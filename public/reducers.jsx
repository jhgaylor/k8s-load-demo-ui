const Reducers = {
  get: () => {
    return (state = {pods: {}}, action) => {
      switch (action.type) {
        case 'ADDED':
          state.pods[action.object.metadata.uid] = action.object
          return state;
        case 'MODIFIED':
          state.pods[action.object.metadata.uid] = action.object
          return state;
        case 'DELETED':
          delete state.pods[action.object.metadata.uid]
          return state;
        case 'ERROR':
          console.error("Got an error event from the api. ?!?", action.object)
          return state;
        default:
          return state
      }
    }
  }
}
