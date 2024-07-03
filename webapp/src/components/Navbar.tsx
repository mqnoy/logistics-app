import { FC } from 'react'
import { NavLink, useNavigate } from 'react-router-dom'
import { cookieUtils, toastUtils } from '../utils'

export const Navbar: FC = () => {
    const navigate = useNavigate()

    const handleLogout = () => {
        cookieUtils.destroyCredentials()
        toastUtils.fireToastSuccess("Logout successfully", {
            onClose: () => {
                navigate(`/login`, { replace: true })
                location.reload()
            }
        })
    }

    return (
        <nav className="navbar is-primary" role="navigation" aria-label="main navigation">
            <div className="navbar-brand">
                <a role="button" className="navbar-burger" aria-label="menu" aria-expanded="false" data-target="navbarBasicExample">
                    <span aria-hidden="true"></span>
                    <span aria-hidden="true"></span>
                    <span aria-hidden="true"></span>
                </a>
            </div>

            <div id="navbarBasicExample" className="navbar-menu">
                <div className="navbar-start">
                    <NavLink className="navbar-item" to={'/'}>
                        Dashboard
                    </NavLink>
                    <NavLink className="navbar-item" to={'/goods'}>
                        Goods
                    </NavLink>
                    <NavLink className="navbar-item" to={'/orders'}>
                        Orders
                    </NavLink>
                </div>

                <div className="navbar-end">
                    <div className="navbar-item">
                        <div className="buttons">
                            <button className="button" onClick={handleLogout}>
                                Logout
                            </button>
                        </div>
                    </div>
                </div>
            </div>
        </nav>
    )
}
