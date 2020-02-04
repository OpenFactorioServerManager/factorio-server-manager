import Swal from 'sweetalert2/dist/sweetalert2';
import withReactContent from 'sweetalert2-react-content';

// define css-classes, that will be appended by sweetalert
const customClassTemp = {
    confirmButton: "swal-btn btn-primary",
    cancelButton: "swal-btn btn-secondary"
}
const customClassDanger = {
    ...customClassTemp,
    confirmButton: "swal-btn btn-danger"
}

// define custom Swals, based on react and with custom designs appended
// USE ONLY THESE, instead of defining it every time again.
const ReactSwal = withReactContent(Swal);
const ReactSwalTemp = ReactSwal.mixin({
    customClass: customClassTemp,
    buttonsStyling: false,

    allowOutsideClick: () => !ReactSwalDanger.isLoading()
});

export const ReactSwalNormal = ReactSwalTemp.mixin({

});

export const ReactSwalDanger = ReactSwalTemp.mixin({
    customClass: customClassDanger
});
