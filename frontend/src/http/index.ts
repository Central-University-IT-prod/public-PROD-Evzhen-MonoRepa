import axios from 'axios'
import { API_URL } from '../constants/constants'

const $api = axios.create({
  withCredentials: false,
  baseURL: API_URL
})

$api.interceptors.request.use((config) => {
  return config
})

$api.interceptors.response.use(
  (config) => {
    return config
  },
  async (error) => {
    console.log(error)
  }
)

const $authApi = axios.create({
  withCredentials: false,
  baseURL: `${API_URL}/authorized`
})

$authApi.interceptors.request.use((config) => {
  config.headers!.Authorization = `Bearer ${localStorage.getItem('token')}`
  return config
})

$authApi.interceptors.response.use(
  (config) => {
    return config
  },
  async (error) => {
    console.log(error)
  }
)

export { $api, $authApi }
