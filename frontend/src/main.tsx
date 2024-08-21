import React, { createContext } from 'react'
import ReactDOM from 'react-dom/client'
import { BrowserRouter, Route, Routes } from 'react-router-dom'
import { createTheme, ThemeProvider } from '@mui/material'
import Store from './store/store.ts'
import { IState } from './models/State/IState.ts'
import ContestList from './components/screens/ContestList/index.tsx'
import CreateContest from './components/screens/CreateContest/index.tsx'
import { ContestInfo } from './components/screens/ContestInfo/ContestInfo.tsx'
import App from './App.tsx'
import { JsCode } from './components/screens/JsCode/JsCode.tsx'

const store = new Store()
export const Context = createContext<IState>({ store })

const theme = createTheme({
  palette: {
    primary: {
      main: '#ffdd2d'
    },
    secondary: {
      main: '#ffdd2d'
    }
  }
})

ReactDOM.createRoot(document.getElementById('root')!).render(
  <React.StrictMode>
    <Context.Provider value={{ store }}>
      <BrowserRouter>
        <ThemeProvider theme={theme}>
          <App>
            <Routes>
              <Route path="/" element={<ContestList />} />
              <Route path="/new-contest" element={<CreateContest />} />
              <Route path="/contest-info/:id" element={<ContestInfo />} />
              <Route path="/js-code" element={<JsCode />} />
            </Routes>
          </App>
        </ThemeProvider>
      </BrowserRouter>
    </Context.Provider>
  </React.StrictMode>
)
