const uriPath = "/api/bus-locations"

// onmessage receives the bounds{} and currentBuses[] on the map. It then gets
// new buses within the bounds and returns the buses to; add, update, and remove
onmessage = (event) => {
	const {
		bounds,
		currentBuses,
	} = event.data;

  // fetch bus locations
	getBusLocations(bounds)
		.then((nextBuses) => {
			// sort buses into; add, update, remove
			postMessage({ buses: sortBuses(currentBuses, nextBuses) })
		})
		.catch((error) => postMessage({ error }))
}

const getBusLocations = (bounds) => {
	const requestUrl =
		`${uriPath}?topLeft=${bounds._sw.lng},${bounds._ne.lat}&bottomRight=${bounds._ne.lng},${bounds._sw.lat}`;

	return new Promise((resolve, reject) => {
		fetch(requestUrl)
			.then(response => response.json())
			.then(response => resolve(response.Buses))
			.catch(error => reject(error))
	})
}

const sortBuses = (currentBuses, updatedBuses) => {
	busesToAdd = [];
	busesToUpdate = [];
	busesToRemove = currentBuses;

	for (
		let updatedBusesIndex = 0;
		updatedBusesIndex < updatedBuses.length;
		updatedBusesIndex++
	) {
		const updatedBus = updatedBuses[updatedBusesIndex];

		busToRemoveIndex = busesToRemove.findIndex(bus => bus.ID == updatedBus.ID);
		if (busToRemoveIndex == -1) {
			// new bus
			busesToAdd.push(updatedBus);
		} else {
			// update bus
			busesToUpdate.push(updatedBus);
			busesToRemove.splice(busToRemoveIndex, 1);
		}
	}

	return {
		add: busesToAdd,
		update: busesToUpdate,
		remove: busesToRemove,
	}
}