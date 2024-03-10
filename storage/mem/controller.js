/**
 *
 * @typedef {Object} PurchaseFlow
 * @property {(reason: string) => Promise<void>} requestclose
 * @property {(receipt: {}) => Promise<void>} receipt
 * @property {(browserId: string, url: string) => Promise<boolean>} launchvalidatedexternalbrowserurl
 * @property {(url: string) => Promise<boolean>} launchexternalbrowserurl
 * @property {(browserId: string) => Promise<string>} getexternalbrowserpath
 * @property {(browserId: string) => Promise<string>} getexternalbrowsername
 * @property {(url: string) => Promise<string>} getdefaultexternalbrowserid
 *
 * @typedef {Object} Engine
 * @property {PurchaseFlow} purchaseflow
 *
 * @typedef {Object} Offer
 * @property {{
 *    displayName: string
 * }} user
 * @property {{
 *   id: string
 *   name: string
 *   price: number
 *   imageUrl: string
 *   type: string
 *  }} offer
 */

function getCookie(name) {
  var value = "; " + document.cookie;
  var parts = value.split("; " + name + "=");
  if (parts.length === 2) return parts.pop().split(";").shift();
}

/**
 * @param {Engine} engine
 * @returns
 */
const main = async (engine) => {
  const close = document.getElementById("close");
  const purchase = document.getElementById("purchaseOfferButton");
  const pf = engine.purchaseflow ? engine.purchaseflow : null;
  if (!pf) return;

  const offerId = new URLSearchParams(window.location.search).get("offers");
  const offerResponse = await axios.get(
    `http://127.0.0.1:3000/purchase/offer?offerId=${offerId}`,
    {
      headers: {
        Authorization: getCookie("EPIC_BEARER_TOKEN"),
      },
    }
  );
  if (offerResponse.status !== 200) return pf.requestclose("LoadFailure");
  const offer = offerResponse.data;

  const image = document.getElementById("orderImage");
  image && (image.style.backgroundImage = `url(${offer.offer.imageUrl})`);

  const title = document.getElementById("orderName");
  title && (title.innerText = offer.offer.name);

  const price = document.getElementById("orderPrice");
  price && (price.innerText = "$" + offer.offer.price);

  const totalPrice = document.getElementById("orderTotalPrice");
  totalPrice && (totalPrice.innerText = "$" + offer.offer.price);

  const SubtotalPrice = document.getElementById("orderSubtotalPrice");
  SubtotalPrice && (SubtotalPrice.innerText = "$" + offer.offer.price);

  const displayName = document.getElementById("displayName");
  displayName && (displayName.innerText = offer.user.displayName);

  close && close.addEventListener("click", () => pf.requestclose("Escape"));
  purchase && purchase.addEventListener("click", () => buy(pf, offer));
};

/**
 * @param {PurchaseFlow} pf
 * @param {Offer} offer
 */
const buy = async (pf, offer) => {
  const purchase = await axios.post(
    `http://127.0.0.1:3000/purchase/offer`,
    {
      offerId: offer.offer.id,
      type: offer.offer.type,
    },
    {
      headers: {
        Authorization: getCookie("EPIC_BEARER_TOKEN"),
      },
    }
  );
  if (purchase.status !== 200) return pf.requestclose("PurchaseFailure");
  await pf.receipt(purchase.data.receipt);
  await pf.requestclose("WasSuccessful");
};
