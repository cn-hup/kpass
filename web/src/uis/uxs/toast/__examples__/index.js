import { toast } from '../toast.manager'

toast.success({
  message: 'Lorem ipsum',
  description: 'Lorem ipsum dolor sit amet, consectetur adipiscing elit.',
  duration: 500
})

toast.error({
  message: 'Lorem ipsum',
  duration: 500
})

toast.info({
  message: '简体中文测试',
  duration: 500
})

setInterval(() => {
  toast.success({
    message: 'Lorem ipsum',
    description: 'Lorem ipsum dolor sit amet, consectetur adipiscing elit.'
  })
}, 7000)
