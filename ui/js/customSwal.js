import Swal from 'sweetalert2/dist/sweetalert2';
import withReactContent from 'sweetalert2-react-content';

const ReactSwal = withReactContent(Swal);
const ReactSwalTemp = ReactSwal.mixin({
    confirmButtonClass: "swal-btn btn-primary",
    cancelButtonClass: "swal-btn btn-secondary",
    customClass: "swal-design",
    buttonsStyling: false,

    allowOutsideClick: () => !ReactSwalDanger.isLoading()
});

export const ReactSwalNormal = ReactSwalTemp.mixin({

});

export const ReactSwalDanger = ReactSwalTemp.mixin({
    confirmButtonClass: "swal-btn btn-danger"
});
