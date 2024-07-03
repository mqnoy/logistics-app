import { FC } from "react"
import { Goods } from "../types"

type GoodDetailProps = {
    data?: Goods
}

export const GoodDetail: FC<GoodDetailProps> = ({ data }) => {
    if (!data) {
        return <>no data</>
    }

    return (
        <div className="">
            <div className="field">
                <label htmlFor="" className="label">
                    Code
                </label>
                <div className="control">
                    <div
                        className="input"
                    >{data.code}</div>
                </div>
            </div>
            <div className="field">
                <label htmlFor="" className="label">
                    Name
                </label>
                <div className="control">
                    <div
                        className="input"
                    >{data.name}</div>
                </div>
            </div>
            <div className="field">
                <label htmlFor="" className="label">
                    Description
                </label>
                <div className="control">
                    <div
                        className="input"
                    >{data.description}</div>
                </div>
            </div>
            <div className="field">
                <label htmlFor="" className="label">
                    Active
                </label>
                <div className="control">
                    <div
                        className="input"
                    >{data.is_active ? 'active' : 'inactive'}</div>
                </div>
            </div>
        </div>
    )
}
