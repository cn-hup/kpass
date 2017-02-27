import React, { Component, PropTypes } from 'react'
import { Link } from 'react-router'
import { Translate } from 'react-redux-i18n'
import cx from 'classnames'

import { Icon } from 'uis'
import { isPublicTeam, isOwner } from 'utils'
import { ENTRIES_PATH, MEMBERS_PATH, SETTINGS_PATH } from '../../index'

import './workspace-sidebar.view.styl'

export class WorkspaceSidebar extends Component {

  static propTypes = {
    className: PropTypes.string,
    team: PropTypes.object,
    userMe: PropTypes.object,
    basePath: PropTypes.string
  }

  getRootClassnames () {
    return cx(
      'workspaceSidebarView',
      this.props.className
    )
  }

  renderSecretNav () {
    const { basePath } = this.props

    return (
      <div className={'navGroup'}>
        <div className={'navGroupTitle'}>
          <Translate value={'pageType.secret'} />
        </div>
        <div className={'navGroupList'}>
          <Link
            to={`${basePath}/${ENTRIES_PATH}`}
            className={'navItem'}
            activeClassName={'isActive'}
          >
            <Icon name={'lock'} />
            <Translate value={'entries.title'} />
          </Link>
        </div>
      </div>
    )
  }

  renderTeamSettingsNav () {
    const { team, userMe, basePath } = this.props

    if (!isOwner(team, userMe)) {
      return null
    }

    return (
      <Link
        to={`${basePath}/${SETTINGS_PATH}`}
        className={'navItem'}
        activeClassName={'isActive'}
      >
        <Icon name={'cog'} />
        <Translate value={'teamSettings.title'} />
      </Link>
    )
  }

  renderTeamNav () {
    const { team, basePath } = this.props

    if (!isPublicTeam(team)) {
      return null
    }

    return (
      <div className={'navGroup'}>
        <div className={'navGroupTitle'}>
          <Translate value={'pageType.team'} />
        </div>
        <div className={'navGroupList'}>
          <Link
            to={`${basePath}/${MEMBERS_PATH}`}
            className={'navItem'}
            activeClassName={'isActive'}
          >
            <Icon name={'users'} />
            <Translate value={'teamMembers.title'} />
          </Link>
          {this.renderTeamSettingsNav()}
        </div>
      </div>
    )
  }

  render () {
    return (
      <div className={this.getRootClassnames()}>
        <div className={'sidebarNav'}>
          {this.renderSecretNav()}
          {this.renderTeamNav()}
        </div>
      </div>
    )
  }

}
