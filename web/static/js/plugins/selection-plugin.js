/**
 * Created by Teocci.
 * Author: teocci@yandex.com on 2023-Sep-03
 */
export default class CheckboxSelectionPlugin {
    /** @type {gridjs.Grid} */
    grid

    constructor(grid) {
        this.grid = grid
    }

    onInit() {
        // Add a column of checkboxes to the grid.
        this.grid.addColumn({
            id: 'select',
            name: 'Select',
            type: 'checkbox',
            width: 40,
            headerTemplate: () => {
                return `<input type="checkbox" id="selectAllCheckbox">`
            },
            cellTemplate(cell, row) {
                return `<input type="checkbox" class="singleCheckbox" ${cell.checked ? 'checked' : ''}>`
            },
        })

        // Add an event listener to the header checkbox.
        const selectAllCheckbox = document.getElementById('selectAllCheckbox')
        selectAllCheckbox.onchange = () => {
            // Check or uncheck all the rows depending on the state of the checkbox.
            this.grid.selectAllRows(event.target.checked)
        }
    }
}