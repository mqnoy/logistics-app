import { FC } from 'react';
import Select, { InputActionMeta } from 'react-select';

type Options = {
    value: number,
    label: string,
}

type DropdownSearchProps<D> = {
    name: string
    isRequired?: boolean
    items: D[];
    renderItem: (item: D) => string
    onSelected: (item: D) => void
    onSearch: (s: string) => void
};

// eslint-disable-next-line @typescript-eslint/no-explicit-any
export const DropdownSearch: FC<DropdownSearchProps<any>> = ({
    name,
    isRequired,
    items,
    renderItem,
    onSelected,
    onSearch
}) => {
    const options = items.map((item, index) => ({
        value: index,
        label: renderItem(item),
    }));

    const handleSelected = (opts: Options) => {
        // fire callback
        const found = items[opts.value]
        onSelected(found)
    };

    const handleInputChange = (newValue: string, _actionMeta: InputActionMeta) => {
        // fire callback
        onSearch(newValue)
    }
    return (
        <div className="container">
            <div className="field">
                <div className="control has-icons-left">
                    <Select
                        name={name}
                        required={isRequired}
                        options={options}
                        className="basic-single"
                        classNamePrefix="select"
                        isClearable={true}
                        isSearchable={true}
                        placeholder="Search..."
                        onChange={(selectedOption) => {
                            if (selectedOption) {
                                handleSelected(selectedOption)
                            }
                        }}
                        onInputChange={handleInputChange}
                    />
                </div>
            </div>
        </div>
    );
};
