import { FC } from 'react'
import { YellowButton } from '../../shared/YellowButton/YellowButton'

export const JsCode: FC = () => {
  const urlParams = new URLSearchParams(window.location.search)
  const code = urlParams.get('code')
  console.log(code)
  return (
    <div style={{ marginLeft: '47vw', marginTop: '35vh' }}>
      <YellowButton>
        <a href={`teamder://helloworld/invite/${code}`}>BEBRA</a>
      </YellowButton>
    </div>
  )
}
