package vars

var (
	CmdDump       = []string{"sh", "-c", "jmap -dump,format=b,file=/tmp/${HOSTNAME}.hprof 1"}
	CmdGetPodIP   = []string{"sh", "-c", "echo ${KUBERNETES_POD_IP}"}
	CmdUploadDump = []string{"sh", "-c", "/usr/bin/oomdump -c 应用触发DUMP -f /tmp/${HOSTNAME}.hprof"}
)
