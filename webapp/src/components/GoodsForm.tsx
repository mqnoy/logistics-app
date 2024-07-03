import { FC, useEffect } from "react"
import { SubmitHandler, useForm } from "react-hook-form"
import { usePostCreateGoodMutation } from "../api"
import { rtkUtils } from "../utils"
import { ErrorApiResponse, ModalActionGoods } from "../types"


type GoodsFromProps = {
    action?: ModalActionGoods
    actionIsDone: (isDone: boolean, error?: ErrorApiResponse) => void
}
type IFormGoods = {
    code: string,
    name: string,
    description: string
    is_active?: boolean
}

export const GoodsForm: FC<GoodsFromProps> = ({ action, actionIsDone }) => {
    const { handleSubmit, watch, setValue } = useForm<IFormGoods>({
        defaultValues: {
            code: '',
            description: '',
            name: ''
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

    const onsubmit: SubmitHandler<IFormGoods> = (data: IFormGoods) => {
        console.debug('onsubmited', data)
        if (action === "create") {
            actionCreate(data)
        }
    }

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