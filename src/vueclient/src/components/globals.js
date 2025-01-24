import moment from 'moment';
import 'moment/dist/locale/ru';
import linkifyHtml from 'linkify-html';

export const glob = {

  fmtDate: (date) => {
    moment.locale('ru');
    return moment(date).fromNow();
  },

  linkify: (src) => {
    const opt = {
      target: {
        url: "_blank",
        email: null,
      },
    };
    return linkifyHtml(src, opt);
  }

};
