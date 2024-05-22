import http from '../utils/http.js'

export const login = (params) => http.post(`/login`, params)