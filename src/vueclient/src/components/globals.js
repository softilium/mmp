import moment from 'moment';
import 'moment/dist/locale/ru';

export const glob = {

  fmtDate: (date) => {
    moment.locale('ru');
    return moment(date).fromNow();
  }

};
