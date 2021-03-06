import { createAction, handleActions } from 'redux-actions'

const initialState = {
  teamId: null
}

export const updateCurrentTeamAction = createAction('UPDATE_CURRENT_TEAM')
export const updateCurrentTeamSuccessAction = createAction('UPDATE_CURRENT_TEAM_SUCCESS')
export const updateCurrentTeamFailureAction = createAction('UPDATE_CURRENT_TEAM_FAILURE')

export const mountCurrentTeamAction = createAction('MOUNT_CURRENT_TEAM')
export const unmountCurrentTeamAction = createAction('UNMOUNT_CURRENT_TEAM')
export const setCurrentTeamAction = createAction('SET_CURRENT_TEAM')

export const currentTeamReducer = handleActions({

  [`${setCurrentTeamAction}`]: (state, action) => {
    if (!action.payload) {
      return state
    }

    const { teamId } = action.payload

    if (typeof teamId === 'undefined') {
      return state
    }

    return {
      ...state,
      teamId
    }
  }

}, initialState)
