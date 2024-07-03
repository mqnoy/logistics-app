import { fetchBaseQuery } from '@reduxjs/toolkit/query'
import { EnvConfig } from '../config'
import { cookieUtils } from '../utils'

export const baseQuery = fetchBaseQuery({
    baseUrl: EnvConfig.apiBaseURL,
    prepareHeaders: (headers) => {
        const authCred = cookieUtils.getCredentials()
        if (authCred) {
            headers.set('Authorization', `Bearer ${authCred.access_token}`)
        }

        headers.set('Access-Control-Allow-Origin', '*')
        return headers
    },
})
