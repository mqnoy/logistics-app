import { createBrowserRouter, redirect } from 'react-router-dom'
import LoginPage from '../pages/Login'
import { DashboardPage } from '../pages'
import { cookieUtils } from '../utils'
import RegisterPage from '../pages/Register'

const privateLoader = () => {
    const authCred = cookieUtils.getCredentials()
    if (!authCred) {
        return redirect('/login')
    }

    return null
}

const router = createBrowserRouter([
    {
        path: '/',
        element: <DashboardPage />,
        loader: privateLoader,
    },
    {
        path: '/login',
        element: <LoginPage />,
        loader: () => {
            const authCred = cookieUtils.getCredentials()
            if (authCred) {
                return redirect('/')
            }
            return null
        },
    },
    {
        path: '/register',
        element: <RegisterPage />,
        loader: () => {
            const authCred = cookieUtils.getCredentials()
            if (authCred) {
                return redirect('/')
            }
            return null
        },
    },
])

export default router
