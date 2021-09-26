import { IndividualSimResult } from '../core/api/api';
import { TypedEvent } from '../core/typed_event';

import { DpsHistogram } from './dps_histogram';

declare var Chart: any;

const urlParams = new URLSearchParams(window.location.search);
if (urlParams.has('mainBgColor')) {
	document.body.style.setProperty('--main-bg-color', urlParams.get('mainBgColor')!);
}
if (urlParams.has('mainTextColor')) {
	document.body.style.setProperty('--main-text-color', urlParams.get('mainTextColor')!);
}

const colorSettings = {
	mainTextColor: document.body.style.getPropertyValue('--main-text-color'),
};

Chart.defaults.color = colorSettings.mainTextColor;

const layoutHTML = `
<div class="dr-root">
	<div class="dps-histogram">
	</div>
</div>
`;

const resultsEmitter = new TypedEvent<IndividualSimResult | null>();
window.addEventListener('message', event => {
	// Null indicates pending results
	const data: IndividualSimResult | null = event.data;

	resultsEmitter.emit(event.data);
});

document.body.innerHTML = layoutHTML;
const dpsHistogram = new DpsHistogram(document.body.getElementsByClassName('dps-histogram')[0] as HTMLElement, resultsEmitter, colorSettings);