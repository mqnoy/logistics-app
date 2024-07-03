import Cookies from 'js-cookie'

export const CookieKeys = {
    access_token: 'access_token',
    refresh_token: 'refresh_token',
    user: 'user',
}

interface AuthCredentials {
    access_token: string
    refresh_token: string
    user: unknown
}

export function setCredentials(creds: AuthCredentials) {
    Cookies.set(CookieKeys.access_token, creds.access_token)
    Cookies.set(CookieKeys.refresh_token, creds.refresh_token)
    Cookies.set(CookieKeys.user, JSON.stringify(creds.user))
}

export function getCredentials(): AuthCredentials | null {
    let cred: AuthCredentials | null = null
    const access_token = Cookies.get(CookieKeys.access_token)
    const refresh_token = Cookies.get(CookieKeys.refresh_token)
    const user = Cookies.get(CookieKeys.user)
    if (access_token && refresh_token && refresh_token && user) {
        cred = {
            access_token,
            refresh_token,
            user: JSON.parse(user),
        } as AuthCredentials
    }

    return cred
}

export function destroyCredentials() {
    Cookies.remove(CookieKeys.access_token)
    Cookies.remove(CookieKeys.refresh_token)
    Cookies.remove(CookieKeys.user)
}
