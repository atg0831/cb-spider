// Cloud Driver Interface of CB-Spider.
// The CB-Spider is a sub-Framework of the Cloud-Barista Multi-Cloud Project.
// The CB-Spider Mission is to connect all the clouds with a single interface.
//
//      * Cloud-Barista: https://github.com/cloud-barista
//
// This is Mock Driver.
//
// by CB-Spider Team, 2020.10.

package resources

import (
	"fmt"
	"time"

	cblog "github.com/cloud-barista/cb-log"
	idrv "github.com/cloud-barista/cb-spider/cloud-control-manager/cloud-driver/interfaces"
	irs "github.com/cloud-barista/cb-spider/cloud-control-manager/cloud-driver/interfaces/resources"
	_ "github.com/sirupsen/logrus"
)

var vmInfoMap map[string][]*irs.VMInfo
var vmStatusInfoMap map[string][]*irs.VMStatusInfo

type MockVMHandler struct {
	Region   idrv.RegionInfo
	MockName string
}

func init() {
	vmInfoMap = make(map[string][]*irs.VMInfo)
	vmStatusInfoMap = make(map[string][]*irs.VMStatusInfo)
}

func (vmHandler *MockVMHandler) StartVM(vmReqInfo irs.VMReqInfo) (irs.VMInfo, error) {
	cblogger := cblog.GetLogger("CB-SPIDER")
	cblogger.Info("Mock Driver: called StartVM()!")

	mockName := vmHandler.MockName
	vmReqInfo.IId.SystemId = vmReqInfo.IId.NameId

	// image validation
	imageHandler := MockImageHandler{mockName}
	validatedImgInfo, err := imageHandler.GetImage(vmReqInfo.ImageIID)
	if err != nil {
		cblogger.Error(err)
		return irs.VMInfo{}, err
	}

	// spec validation
	vmSpecHandler := MockVMSpecHandler{mockName}
	validatedSpecInfo, err := vmSpecHandler.GetVMSpec(vmHandler.Region.Region, vmReqInfo.VMSpecName)
	if err != nil {
		cblogger.Error(err)
		return irs.VMInfo{}, err
	}

	// vpc validation
	vpcHandler := MockVPCHandler{mockName}
	validatedVPCInfo, err := vpcHandler.GetVPC(vmReqInfo.VpcIID)
	if err != nil {
		cblogger.Error(err)
		return irs.VMInfo{}, err
	}

	// subnet validation
	var validatedSubnetInfo *irs.SubnetInfo = nil
	for _, info := range validatedVPCInfo.SubnetInfoList {
		if info.IId.NameId == vmReqInfo.SubnetIID.NameId {
			validatedSubnetInfo = &info
		}
	}
	if validatedSubnetInfo == nil {
		return irs.VMInfo{}, fmt.Errorf("%s subnet iid does not exist!!", vmReqInfo.SubnetIID.NameId)
	}

	// sg validation
	securityHandler := MockSecurityHandler{mockName}
	sgInfoList, err := securityHandler.ListSecurity()
	if err != nil {
		cblogger.Error(err)
		return irs.VMInfo{}, err
	}
	validatedSgIIDs := []irs.IID{}
	for _, info1 := range vmReqInfo.SecurityGroupIIDs {
		flg := false
		for _, info2 := range sgInfoList {
			if (*info2).IId.NameId == info1.NameId {
				validatedSgIIDs = append(validatedSgIIDs, info2.IId)
				flg = true
				break
			}
		}
		if !flg {
			return irs.VMInfo{}, fmt.Errorf("%s security group iid does not exist!!", info1.NameId)
		}
	}

	// keypair validation
	keyPairHandler := MockKeyPairHandler{mockName}
	validatedKeyPairInfo, err := keyPairHandler.GetKey(vmReqInfo.KeyPairIID)
	if err != nil {
		cblogger.Error(err)
		return irs.VMInfo{}, err
	}

	// vm creation
	vmInfo := irs.VMInfo{
		IId:       vmReqInfo.IId,
		StartTime: time.Now(),

		Region:            irs.RegionInfo{vmHandler.Region.Region, vmHandler.Region.Zone},
		ImageIId:          validatedImgInfo.IId,
		VMSpecName:        validatedSpecInfo.Name,
		VpcIID:            validatedVPCInfo.IId,
		SubnetIID:         validatedSubnetInfo.IId,
		SecurityGroupIIds: validatedSgIIDs,

		KeyPairIId: validatedKeyPairInfo.IId,

		VMUserId:     vmReqInfo.VMUserId,
		VMUserPasswd: vmReqInfo.VMUserPasswd,

		NetworkInterface: "",
		PublicIP:         "",
		PublicDNS:        "",
		PrivateIP:        "",
		PrivateDNS:       "",

		VMBootDisk:  "",
		VMBlockDisk: "",

		KeyValueList: nil,
	}

	infoList, _ := vmInfoMap[mockName]
	infoList = append(infoList, &vmInfo)
	vmInfoMap[mockName] = infoList

	// vm status creation
	vmStatusInfo := irs.VMStatusInfo{vmReqInfo.IId, "Running"}

	statusInfoList, _ := vmStatusInfoMap[mockName]
	statusInfoList = append(statusInfoList, &vmStatusInfo)
	vmStatusInfoMap[mockName] = statusInfoList

	return vmInfo, nil
}

func (vmHandler *MockVMHandler) SuspendVM(iid irs.IID) (irs.VMStatus, error) {
	cblogger := cblog.GetLogger("CB-SPIDER")
	cblogger.Info("Mock Driver: called SuspendVM()!")

	mockName := vmHandler.MockName

	statusInfoList, ok := vmStatusInfoMap[mockName]
	if !ok {
		return "", fmt.Errorf("%s vm status does not exist!!", mockName)
	}

	var validatedStatusInfo *irs.VMStatusInfo = nil
	for _, info := range statusInfoList {
		if (*info).IId.NameId == iid.NameId {
			validatedStatusInfo = info
		}
	}
	if validatedStatusInfo == nil {
		return "", fmt.Errorf("%s status iid does not exist!!", iid.NameId)
	}

	validatedStatusInfo.VmStatus = "Suspended"
	return "Suspending", nil
}

func (vmHandler *MockVMHandler) ResumeVM(iid irs.IID) (irs.VMStatus, error) {
	cblogger := cblog.GetLogger("CB-SPIDER")
	cblogger.Info("Mock Driver: called ResumeVM()!")

	mockName := vmHandler.MockName

	statusInfoList, ok := vmStatusInfoMap[mockName]
	if !ok {
		return "", fmt.Errorf("%s vm status does not exist!!", mockName)
	}

	var validatedStatusInfo *irs.VMStatusInfo = nil
	for _, info := range statusInfoList {
		if (*info).IId.NameId == iid.NameId {
			validatedStatusInfo = info
		}
	}
	if validatedStatusInfo == nil {
		return "", fmt.Errorf("%s status iid does not exist!!", iid.NameId)
	}

	validatedStatusInfo.VmStatus = "Running"
	return "Resuming", nil
}

func (vmHandler *MockVMHandler) RebootVM(iid irs.IID) (irs.VMStatus, error) {
	cblogger := cblog.GetLogger("CB-SPIDER")
	cblogger.Info("Mock Driver: called RebootVM()!")

	mockName := vmHandler.MockName

	statusInfoList, ok := vmStatusInfoMap[mockName]
	if !ok {
		return "", fmt.Errorf("%s vm status does not exist!!", mockName)
	}

	var validatedStatusInfo *irs.VMStatusInfo = nil
	for _, info := range statusInfoList {
		if (*info).IId.NameId == iid.NameId {
			validatedStatusInfo = info
		}
	}
	if validatedStatusInfo == nil {
		return "", fmt.Errorf("%s status iid does not exist!!", iid.NameId)
	}

	if validatedStatusInfo.VmStatus == "Suspended" {
		return "", fmt.Errorf("reboot not supported in SUSPENDED status")
	}

	validatedStatusInfo.VmStatus = "Running"
	return "Rebooting", nil
}

func (vmHandler *MockVMHandler) TerminateVM(iid irs.IID) (irs.VMStatus, error) {
	cblogger := cblog.GetLogger("CB-SPIDER")
	cblogger.Info("Mock Driver: called TerminateVM()!")

	infoList, err := vmHandler.ListVM()
	if err != nil {
		cblogger.Error(err)
		return "", err
	}

	statusInfoList, err := vmHandler.ListVMStatus()
	if err != nil {
		cblogger.Error(err)
		return "", err
	}

	mockName := vmHandler.MockName
	for idx, info := range infoList {
		if info.IId.NameId == iid.NameId {
			infoList = append(infoList[:idx], infoList[idx+1:]...)
		}
	}
	vmInfoMap[mockName] = infoList

	for idx, info := range statusInfoList {
		if info.IId.NameId == iid.NameId {
			statusInfoList = append(statusInfoList[:idx], statusInfoList[idx+1:]...)
		}
	}
	vmStatusInfoMap[mockName] = statusInfoList

	return "Terminating", nil
}

func (vmHandler *MockVMHandler) ListVMStatus() ([]*irs.VMStatusInfo, error) {
	cblogger := cblog.GetLogger("CB-SPIDER")
	cblogger.Info("Mock Driver: called ListVMStatus()!")

	mockName := vmHandler.MockName
	infoList, ok := vmStatusInfoMap[mockName]
	if !ok {
		return []*irs.VMStatusInfo{}, nil
	}
	// cloning list of VM Status
	resultList := make([]*irs.VMStatusInfo, len(infoList))
	copy(resultList, infoList)
	return resultList, nil
}

func (vmHandler *MockVMHandler) GetVMStatus(iid irs.IID) (irs.VMStatus, error) {
	cblogger := cblog.GetLogger("CB-SPIDER")
	cblogger.Info("Mock Driver: called GetVMStatus()!")

	infoList, err := vmHandler.ListVMStatus()
	if err != nil {
		cblogger.Error(err)
		return "", err
	}

	for _, info := range infoList {
		if (*info).IId.NameId == iid.NameId {
			return (*info).VmStatus, nil
		}
	}

	return "", fmt.Errorf("%s status iid does not exist!!", iid.NameId)
}

func (vmHandler *MockVMHandler) ListVM() ([]*irs.VMInfo, error) {
	cblogger := cblog.GetLogger("CB-SPIDER")
	cblogger.Info("Mock Driver: called ListVM()!")

	mockName := vmHandler.MockName
	infoList, ok := vmInfoMap[mockName]
	if !ok {
		return []*irs.VMInfo{}, nil
	}
	// cloning list of VM
	resultList := make([]*irs.VMInfo, len(infoList))
	copy(resultList, infoList)
	return resultList, nil
}

func (vmHandler *MockVMHandler) GetVM(iid irs.IID) (irs.VMInfo, error) {
	cblogger := cblog.GetLogger("CB-SPIDER")
	cblogger.Info("Mock Driver: called GetVM()!")

	infoList, err := vmHandler.ListVM()
	if err != nil {
		cblogger.Error(err)
		return irs.VMInfo{}, err
	}

	for _, info := range infoList {
		if (*info).IId.NameId == iid.NameId {
			return *info, nil
		}
	}

	return irs.VMInfo{}, fmt.Errorf("%s vm iid does not exist!!", iid.NameId)
}