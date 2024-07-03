import { FC, useEffect, useState } from 'react'
import { SubmitHandler, useForm } from 'react-hook-form'
import { NavLink, useNavigate } from 'react-router-dom'
import { usePostUserLoginMutation } from '../api/authApi'
import { cookieUtils, rtkUtils, toastUtils } from '../utils'

type IFormLogin = {
    email: string
    password: string
}

const LoginPage: FC = () => {
    const navigate = useNavigate()
    const { handleSubmit, watch, setValue } = useForm<IFormLogin>({
        defaultValues: {
            email: '',
            password: '',
        },
    })

    const [postUserLogin, { isError, isLoading, isSuccess, data: dataUserLogin, error: errorResp }] =
        usePostUserLoginMutation()
    const [buttonEnable, setButtonEnable] = useState(true)

    useEffect(() => {
        if (isSuccess) {
            cookieUtils.setCredentials({
                access_token: dataUserLogin.data.access_token,
                refresh_token: dataUserLogin.data.refresh_token,
                user: dataUserLogin.data.user,
            })
            toastUtils.fireToastSuccess("Login successfully", {
                onClose() {
                    navigate('/')
                },
            })

        } else if (isError) {
            if (errorResp) {
                const errorApi = rtkUtils.parseErrorRtk(errorResp);
                toastUtils.fireToastError(errorApi)
            }

        } else if (isLoading) {
            setButtonEnable(false)
        }
    }, [isSuccess, isError, isLoading, buttonEnable])


    const email = watch('email')
    const password = watch('password')

    const onsubmit: SubmitHandler<IFormLogin> = (data: IFormLogin) => {
        console.debug('onsubmited', data)
        postUserLogin(data)
    }

    return (
        <>
            <section className="hero is-fullheight is-full-width">
                <div className="hero-body">
                    <div className="container ">
                        <div className="columns is-centered">
                            <div className="column is-4">
                                <form className="box" onSubmit={handleSubmit(onsubmit)}>
                                    <h1 className="title is-3">Sign In</h1>
                                    <div className="field">
                                        <label htmlFor="" className="label">
                                            Email
                                        </label>
                                        <div className="control">
                                            <input
                                                type="email"
                                                value={email}
                                                onChange={(e) => {
                                                    setValue('email', e.target.value)
                                                }}
                                                placeholder="example@domain.tld"
                                                className="input"
                                                required
                                                autoComplete="email"
                                            />
                                        </div>
                                    </div>
                                    <div className="field">
                                        <label htmlFor="" className="label">
                                            Password
                                        </label>
                                        <div className="control">
                                            <input
                                                type="password"
                                                value={password}
                                                onChange={(e) => {
                                                    setValue('password', e.target.value)
                                                }}
                                                placeholder="****"
                                                className="input"
                                                required
                                            />
                                        </div>
                                    </div>
                                    <div className="field">
                                        <button
                                            disabled={!buttonEnable}
                                            className="button is-success is-fullwidth"
                                            type="submit"
                                        >
                                            Login
                                        </button>
                                    </div>
                                    <div className='has-text-centered'>
                                        Don't have an account yet <NavLink className="" to={'/register'}>register</NavLink>
                                    </div>
                                </form>
                            </div>
                        </div>
                    </div>
                </div>
            </section>
        </>
    )
}

export default LoginPage
