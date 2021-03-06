import { I18n } from 'react-redux-i18n'

export const getEntryCategories = () => {
  const baseCategories = [
    {
      value: 'Login',
      color: '#F58E3D'
    },

    {
      value: 'Network',
      color: '#797EC9'
    },

    {
      value: 'Software License',
      color: '#3DA8F5'
    },

    {
      value: 'Secure Note',
      color: '#75C940'
    },

    {
      value: 'Server',
      color: '#FFE738'
    },

    {
      value: 'Wallet',
      color: '#FF4F3E'
    }
  ]

  // I18n
  const enhancedCategories = baseCategories.map((category) => {
    category.title = I18n.t(`entryCategory.${category.value}`)
    return category
  })

  return enhancedCategories
}

export const getEntryCategoryByValue = (value) => {
  return getEntryCategories().find(
    category => category.value === value
  )
}
