import { createBrowserRouter, redirect } from 'react-router-dom'
import LoginPage from '../pages/Login'
import { DashboardPage } from '../pages'
import { cookieUtils } from '../utils'
import RegisterPage from '../pages/Register'
import GoodsPage from '../pages/GoodsPage'
import { OrderPage } from '../pages/OrderPage'
import { OrderCreatePage } from '../pages/OrderPage/OrderCreatePage'
import { OrderCreatePageV2 } from '../pages/OrderPage/OrderCreatePageV2'

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
    {
        path: '/goods',
        element: <GoodsPage />,
        loader: privateLoader,
    },
    {
        path: '/orders',
        element: <OrderPage />,
        loader: privateLoader,
    },
    {
        path: '/orders/create',
        element: <OrderCreatePage />,
        loader: privateLoader,
    },
    {
        path: '/orders/create/v2',
        element: <OrderCreatePageV2 />,
        loader: privateLoader,
    },
])

export default router
