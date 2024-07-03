import { FC, useEffect } from "react"
import { SubmitHandler, useForm } from "react-hook-form"
import { usePostCreateGoodMutation, usePutUpdateGoodMutation } from "../api"
import { rtkUtils } from "../utils"
import { ErrorApiResponse, Goods, ModalActionGoods } from "../types"


type GoodsFromProps = {
    action?: ModalActionGoods
    actionIsDone: (isDone: boolean, error?: ErrorApiResponse) => void
    dataDetail?: Goods
}
type IFormGoods = {
    code: string,
    name: string,
    description: string
    is_active?: string
    id?: string
}

export const GoodsForm: FC<GoodsFromProps> = ({ action, actionIsDone, dataDetail }) => {
    const { register, reset, handleSubmit, watch, setValue } = useForm<IFormGoods>({
        defaultValues: {
            code: '',
            description: '',
            name: '',
            is_active: 'active'
        },
    })
    const code = watch('code')
    const name = watch('name')
    const description = watch('description')

    const [postCreateGoods, { isError: isErrorCreateGoods, isLoading: isLoadingCreateGoods, isSuccess: isSuccessCreateGoods, error: errorRespCreateGoods }] =
        usePostCreateGoodMutation()
    useEffect(() => {
        if (isSuccessCreateGoods) {
            actionIsDone(true)
        } else if (errorRespCreateGoods) {
            const errorApi = rtkUtils.parseErrorRtk(errorRespCreateGoods);

            actionIsDone(true, errorApi)
        } else if (isLoadingCreateGoods) {
            console.debug('loading..');
            actionIsDone(false)
        }
    }, [isSuccessCreateGoods, isErrorCreateGoods, isLoadingCreateGoods])


    const actionCreate = (data: IFormGoods) => {
        postCreateGoods({
            code: data.code,
            description: data.description,
            name: data.name
        })
    }

    const [putUpdateGoods, { isError: isErrorUpdateGoods, isLoading: isLoadingUpdateGoods, isSuccess: isSuccessUpdateGoods, error: errorRespUpdateGoods }] =
        usePutUpdateGoodMutation()
    useEffect(() => {
        if (isSuccessUpdateGoods) {
            actionIsDone(true)
        } else if (errorRespUpdateGoods) {
            const errorApi = rtkUtils.parseErrorRtk(errorRespUpdateGoods);

            actionIsDone(true, errorApi)
        } else if (isLoadingUpdateGoods) {
            console.debug('loading..');
            actionIsDone(false)
        }
    }, [isSuccessUpdateGoods, isErrorUpdateGoods, isLoadingUpdateGoods])

    const actionUpdate = (data: IFormGoods) => {
        if (data.id && data.is_active) {
            putUpdateGoods({
                id: data.id,
                code: data.code,
                description: data.description,
                name: data.name,
                is_active: data.is_active === "active" ? true : false
            })
        }
    }

    const onsubmit: SubmitHandler<IFormGoods> = (data: IFormGoods) => {
        console.debug('onsubmited', data)
        if (action === "create") {
            actionCreate(data)
        } else if (action === "update") {
            actionUpdate(data)
        }
    }

    useEffect(() => {
        if (dataDetail) {
            setValue('id', dataDetail.id)
            setValue('code', dataDetail.code)
            setValue('name', dataDetail.name)
            setValue('description', dataDetail.description)
            setValue('is_active', dataDetail.is_active ? "active" : "inactive")
        }

        if (action === "create") {
            reset()
        }
    }, [dataDetail, action])


    return (
        <form className="" onSubmit={handleSubmit(onsubmit)}>
            <div className="field">
                <label htmlFor="" className="label">
                    Code
                </label>
                <div className="control">
                    <input
                        type="text"
                        value={code}
                        onChange={(e) => {
                            setValue('code', e.target.value)
                        }}
                        placeholder="KBYD-00"
                        className="input"
                        required
                    />
                </div>
            </div>
            <div className="field">
                <label htmlFor="" className="label">
                    Name
                </label>
                <div className="control">
                    <input
                        type="text"
                        value={name}
                        onChange={(e) => {
                            setValue('name', e.target.value)
                        }}
                        placeholder="Keyboard"
                        className="input"
                        required
                    />
                </div>
            </div>
            <div className="field">
                <label htmlFor="" className="label">
                    Description
                </label>
                <div className="control">
                    <textarea
                        value={description}
                        onChange={(e) => {
                            setValue('description', e.target.value)
                        }}
                        placeholder="keyboard digital alliance"
                        className="input"
                        required
                    />
                </div>
            </div>
            {action === "update" && <div className="field">
                <label htmlFor="" className="label">
                    Active
                </label>
                <div className="control">
                    <label className="radio">
                        <input type="radio"
                            value="active"
                            {...register('is_active', { required: true })}
                        />
                        Active
                    </label>
                    <label className="radio">
                        <input type="radio"
                            value="inactive"
                            {...register('is_active', { required: true })}
                        />
                        Inactive
                    </label>
                </div>
            </div>
            }
            <div className="field">
                <button
                    className="button is-success is-fullwidth"
                    type="submit"
                >
                    Submit
                </button>
            </div>
        </form>
    )
}