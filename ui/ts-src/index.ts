import './styles.css';

import '@shoelace-style/shoelace/dist/themes/light.css';
import '@shoelace-style/shoelace/dist/components/alert/alert.js';
import '@shoelace-style/shoelace/dist/components/button/button.js';
import '@shoelace-style/shoelace/dist/components/divider/divider.js';
import '@shoelace-style/shoelace/dist/components/dropdown/dropdown.js';
import '@shoelace-style/shoelace/dist/components/menu/menu.js';
import '@shoelace-style/shoelace/dist/components/menu-item/menu-item.js';
import '@shoelace-style/shoelace/dist/components/input/input.js';
import '@shoelace-style/shoelace/dist/components/icon/icon.js';
import '@shoelace-style/shoelace/dist/components/tag/tag.js';
import '@shoelace-style/shoelace/dist/components/dialog/dialog.js';
import {setBasePath} from '@shoelace-style/shoelace/dist/utilities/base-path.js';

setBasePath('/vendor/shoelace');

window['isEmptyFilter'] = function (value: string, filterWrapper: string): string {
  if (!value.includes(':') &&
      !value.includes('>') &&
      !value.includes('<') &&
      !value.includes('=') &&
      !value.includes('AND') &&
      !value.includes('OR')) {
    return filterWrapper.replace('{value}', value);
  }

  return value;
}
