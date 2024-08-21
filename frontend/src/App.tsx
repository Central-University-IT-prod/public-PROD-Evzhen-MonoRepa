import { FC, ReactNode, useContext, useLayoutEffect } from 'react'
import { observer } from 'mobx-react-lite'
import { Context } from './main'
import { Auth } from './components/screens/Auth/Auth'
import './styles/global.scss'
import { Loading } from './components/shared/Loading/Loading'
import Header from './components/layouts/Header/Header'
import { useLocation } from 'react-router-dom'

interface IAppProps {
  children: ReactNode
}

const App: FC<IAppProps> = ({ children }) => {
  const { store } = useContext(Context)

  useLayoutEffect(() => {
    if (localStorage.getItem('token')) {
      store.setAuth(true)
    }
  }, [])

  if (store.isLoading) {
    return <Loading />
  }
  if (!store.isAuth && window.location.pathname !== '/js-code') {
    return <Auth />
  }

  return (
    <div className="wrapper">
      <Header />

      <main className="main">{children}</main>
    </div>
  )
}

export default observer(App)
