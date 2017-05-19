import {connect} from 'react-redux';
import {addAccount} from '../action/account';
import {loadAccounts} from '../action/account';


import NewAccount from '../components/CreateAccount/NewAccount'

const mapStateToProps = (state) => {
    return {
        accounts: state.accounts,
    };
};
const mapDispatchToProps = {
    onSubmitAccount: addAccount,
    loadAccounts: loadAccounts
};

export default connect(mapStateToProps, mapDispatchToProps)(NewAccount);