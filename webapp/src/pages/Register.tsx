import { FC, useEffect, useState } from 'react'
import { SubmitHandler, useForm } from 'react-hook-form'
import { NavLink, useNavigate } from 'react-router-dom'
import { usePostUserRegisterMutation } from '../api'
import { rtkUtils, toastUtils } from '../utils'
type IFormRegister = {
    full_name: string
    email: string
    password: string
    passwordMatch: string
}

const RegisterPage: FC = () => {
    const navigate = useNavigate()
    const { handleSubmit, watch, setValue } = useForm<IFormRegister>({
        defaultValues: {
            full_name: '',
            email: '',
            password: '',
            passwordMatch: '',
        },
    })

    const [postUserRegister, { isError, isLoading, isSuccess, error: errorResp }] =
        usePostUserRegisterMutation()
    const [buttonEnable, setButtonEnable] = useState(false)


    useEffect(() => {
        if (isSuccess) {
            toastUtils.fireToastSuccess("Register successfully", {
                onClose() {
                    navigate('/login')
                },
            })
            navigate('/login')
        } else if (isError) {
            if (errorResp) {
                const errorApi = rtkUtils.parseErrorRtk(errorResp);
                toastUtils.fireToastError(errorApi)
            }

        } else if (isLoading) {
            setButtonEnable(false)
        }
    }, [isSuccess, isError, isLoading, buttonEnable])


    const fullName = watch('full_name')
    const email = watch('email')
    const password = watch('password')
    const passwordMatch = watch('passwordMatch')

    const matchPassword = () => {
        if (password !== '' && passwordMatch !== '' && password === passwordMatch) {
            setButtonEnable(true)
        } else {
            setButtonEnable(false)
        }
        console.debug(`matchpassword invoked: password: ${password} , passwordMatch: ${passwordMatch}`);
    }


    const onsubmit: SubmitHandler<IFormRegister> = (data: IFormRegister) => {
        console.debug('onsubmited', data)
        postUserRegister({
            full_name: data.full_name,
            email: data.email,
            password: data.password
        })
    }

    return (
        <>
            <section className="hero is-fullheight is-full-width">
                <div className="hero-body">
                    <div className="container ">
                        <div className="columns is-centered">
                            <div className="column is-4">
                                <form className="box" onSubmit={handleSubmit(onsubmit)}>
                                    <h1 className="title is-3">Register</h1>
                                    <div className="field">
                                        <label htmlFor="" className="label">
                                            Full Name
                                        </label>
                                        <div className="control">
                                            <input
                                                type="text"
                                                value={fullName}
                                                onChange={(e) => {
                                                    setValue('full_name', e.target.value)
                                                }}
                                                placeholder="jhon doe"
                                                className="input"
                                                required
                                                autoComplete="name"
                                            />
                                        </div>
                                    </div>
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
                                                onFocus={() => {
                                                    matchPassword()
                                                }}
                                                onChange={(e) => {
                                                    setValue('password', e.target.value)
                                                }}
                                                onKeyUp={() => {
                                                    matchPassword()
                                                }}
                                                placeholder="****"
                                                className="input"
                                                required
                                            />
                                        </div>
                                    </div>
                                    <div className="field">
                                        <label htmlFor="" className="label">
                                            Repeat password
                                        </label>
                                        <div className="control">
                                            <input
                                                type="password"
                                                value={passwordMatch}
                                                onFocus={() => {
                                                    matchPassword()
                                                }}
                                                onChange={(e) => {
                                                    setValue('passwordMatch', e.target.value)
                                                }}
                                                onKeyUp={() => {
                                                    matchPassword()
                                                }}
                                                placeholder="****"
                                                className="input"
                                                required
                                            />
                                        </div>
                                    </div>
                                    <div className="field">
                                        <button
                                            className="button is-success is-fullwidth"
                                            disabled={!buttonEnable}
                                            type="submit"
                                        >
                                            Register
                                        </button>
                                    </div>
                                    <div className='has-text-centered'>
                                        Already have account? <NavLink className="" to={'/login'}>login</NavLink>
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

export default RegisterPage
