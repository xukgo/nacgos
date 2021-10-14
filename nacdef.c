#include <stdlib.h>
#include "nacdef.h"

void nacdef_doUpdateEvent(ConfigUpdateEvent evt, char *group, char *dataId, char *data){
	evt(group,dataId,data);
}
